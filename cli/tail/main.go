package main

import (
	"bufio"
	"cli/api"
	"cli/constants"
	"cli/models"
	"cli/tail/utils"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strings"
	"time"
)

var batchedLogs []models.LogEntry
var backendDisabled bool = false

func argParser() (int, int, string) {
	poll := flag.Int("poll", constants.PollIntervalLimit, "Default polling time in milliseconds")
	batchSize := flag.Int("batch", constants.BatchLimit, "Default batch size")
	tag := flag.String("tag", "", "Tag for the session")

	flag.Parse() // MUST PARSE FLAGS

	return *poll, *batchSize, *tag
}

func main() {
	pollInterval, batchSize, tag := argParser()
	fmt.Println("Poll:", pollInterval, "Batch:", batchSize)

	sessionId, err := api.SetupSession()
	if err != nil {
		log.Fatal(err)
	}

	filePath := "sample.log"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(filePath); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Watching:", filePath)

	var lastSize int = 0

	for {
		select {

		case event := <-watcher.Events:
			if event.Op&fsnotify.Write != fsnotify.Write {
				continue
			}
			fmt.Println("New write found")
			f, _ := os.Open(filePath)
			f.Seek(int64(lastSize), 0)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "EOF" {
					api.BatchUpload(batchedLogs, sessionId, tag, 3, &backendDisabled)
					os.Exit(0)
				}
				if len(batchedLogs) == batchSize {
					api.BatchUpload(batchedLogs, sessionId, tag, 3, &backendDisabled)
					batchedLogs = batchedLogs[:0]
				}
				const start = "<uplog>"
				const end = "</uplog>"

				if strings.HasPrefix(line, start) && strings.HasSuffix(line, end) {
					content := line[len(start) : len(line)-len(end)]
					entry := models.LogEntry{
						Message:   line,
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Level:     utils.DetectLevel(line),
					}
					batchedLogs = append(batchedLogs, entry)

					fmt.Println("LOG CONTENT:", content)
				}
				lastSize = lastSize + len(line)
			}

		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
