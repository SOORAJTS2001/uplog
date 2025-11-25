package worker


import (
	"os"
	"io"
	"bufio"
	"sync"
	"fmt"
	"time"
	"context"
	"bytes"
	"net/http"
	"database/sql"
	"encoding/json"
	"uplog/models"
	"uplog/config"
	"uplog/utils"
	"uplog/db"
)

func CopyStreamsToFile(sql_object *sql.DB,stdout io.ReadCloser, stderr io.ReadCloser, outFile *os.File, sessionID string) {
	// We'll prefix lines with nothing, but update db for line_count and size_bytes
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	var wg sync.WaitGroup
	wg.Add(2)

	writeLines := func(scanner *bufio.Scanner) {
		defer wg.Done()
		for scanner.Scan() {
			line := scanner.Text()
			// write line + newline
			fmt.Println(line) // <-- show to terminal
			entry := models.LogEntry{
				Message:   line,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Level:     utils.DetectLevel(line),
			}

			// Serialize to JSONL or a custom format
			b, _ := json.Marshal(entry)
			outFile.Write(b)
			outFile.WriteString("\n")

			db.IncrementSessionStats(sql_object,sessionID, 1, int64(len(b)))
			// flush to disk for durability
			_ = outFile.Sync()
		}
	}
	go writeLines(stdoutScanner)
	go writeLines(stderrScanner)
	wg.Wait()
}


func UploaderLoop(ctx context.Context, path string, sessionID string,pollInterval time.Duration,tag*string, doneWriting <-chan struct{}) error {
	var offset int64 = 0
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		// uploader should stop only when writer has completely finished writing logs
		case <-doneWriting:
			// final drain
			if err := UploadNewChunks(path, &offset, sessionID,tag); err != nil {
				return err
			}
			return nil

		// regular timed uploads
		case <-ticker.C:
			if err := UploadNewChunks(path, &offset, sessionID,tag); err != nil {
				fmt.Printf("upload error (will retry): %v\n", err)
			}

		// context cancel should NOT stop uploader
		// it only stops polling but uploader continues until doneWriting closes
		case <-ctx.Done():
			// do nothing — wait for doneWriting

		}
	}
}


func UploadNewChunks(path string, offset *int64, sessionID string,tag *string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	size := info.Size()
	if *offset >= size {
		return nil // nothing new
	}

	// Seek to last processed offset
	_, err = f.Seek(*offset, io.SeekStart)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	batch := make([]models.LogEntry, 0, config.BatchLimit)

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		*offset += int64(len(lineBytes)) + 1 // +1 for newline

		var entry models.LogEntry
		if err := json.Unmarshal(lineBytes, &entry); err != nil {
			continue
		}

		batch = append(batch, entry)

		if len(batch) >= config.BatchLimit {
			if err := SendBatch(sessionID, batch,tag); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	// Send leftover batch
	if len(batch) > 0 {
		return SendBatch(sessionID, batch, tag)
	}

	return scanner.Err()
}

// sendBatch sends a chunk to the backend
func SendBatch(sessionID string, batch []models.LogEntry, tag *string) error {
	body, _ := json.Marshal(batch)
	req, err := http.NewRequest("POST",
		config.BackendUploadEndpoint+"?session_id="+sessionID+"&"+"tag="+*tag,
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if key := os.Getenv("UPLOG_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %d %s", resp.StatusCode, string(data))
	}

	return nil
}
