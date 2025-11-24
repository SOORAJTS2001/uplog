package main

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	baseDir                      = ".uplog"
	configDir                    = "config"
	tmpDir                       = "tmp"
	sqliteFileName               = "db.sqlite"
	domain                       = "https://logs.uplog.com"
	logsDomain                   = domain + "/session="
	backendDomain                = "http://127.0.0.1:8000"
	backendUploadEndpoint        = backendDomain + "/session/upload"
	backendSessionCreateEndpoint = backendDomain + "/session/create"
	chunkSize                    = 32 * 1024 // 32 KB
	batchLimit                   = 200
	pollIntervalLimit            = 200 * time.Millisecond
)

var db *sql.DB
var pollInterval = 200 * time.Millisecond
var respJSON SessionCreateResponse
var tag *string


type Session struct {
	SessionID  string
	CreatedAt  time.Time
	ExpiredAt  sql.NullTime
	LineCount  int64
	SizeBytes  int64
	IsUploaded bool
	Mode       string
	Tag        string
}

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
}

type SessionCreateResponse struct {
    SessionID string `json:"session_id"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	home := userHomeDir()
	if home == "" {
		fmt.Println("cannot get user home directory")
		os.Exit(1)
	}
	setupDirectories(home)
	initDB(home)

	cmd := os.Args[1]

	switch cmd {

	case "run":
		runCmdWithFlags(home, os.Args[2:])

	case "list":
		listCmd(home)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("usage: uplog delete <session_id>")
			os.Exit(1)
		}
		deleteCmd(home, os.Args[2])

	default:
		usage()
	}
}

func usage() {
	fmt.Println("uplog - simple CLI")
	fmt.Println("Usage:")
	fmt.Println("  uplog run <cmd> [args...]    Run a command and upload logs")
	fmt.Println("  uplog list                   List uploaded sessions")
}

// -------------------- FS + Init --------------------

func userHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// fallback
	dir, _ := os.UserHomeDir()
	return dir
}

func setupDirectories(home string) {
	base := filepath.Join(home, baseDir)
	_ = os.MkdirAll(filepath.Join(base, configDir), 0o700)
	_ = os.MkdirAll(filepath.Join(base, tmpDir), 0o700)
}

func dbPath(home string) string {
	return filepath.Join(home, baseDir, configDir, sqliteFileName)
}

func initDB(home string) {
	var err error
	db, err = sql.Open("sqlite3", dbPath(home))
	if err != nil {
		fmt.Printf("failed to open sqlite: %v\n", err)
		os.Exit(1)
	}
	// busy timeout
	_, _ = db.Exec("PRAGMA busy_timeout = 5000;")

	create := `
	CREATE TABLE IF NOT EXISTS sessions (
		session_id   TEXT PRIMARY KEY,
		created_at   DATETIME NOT NULL,
		expired_at   DATETIME,
		line_count   INTEGER DEFAULT 0,
		size_bytes   INTEGER DEFAULT 0,
		is_uploaded  INTEGER DEFAULT 0,
		mode         TEXT NOT NULL,
		tag			 TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at);
	CREATE INDEX IF NOT EXISTS idx_sessions_expired_at ON sessions(expired_at);
	`
	_, err = db.Exec(create)
	if err != nil {
		fmt.Printf("failed to create tables: %v\n", err)
		os.Exit(1)
	}
}

// -------------------- Commands --------------------

func runCmd(home, command string, args []string) {
	// 1. Get session ID from backend (or local generation if backend not available)
	sessionID, err := requestSessionIDFromBackend()
	doneWriting := make(chan struct{})
	if err != nil {
		// fallback to local uuid but still continue
		fmt.Printf("warning: backend session request failed: %v. Using local UUID.\n", err)
		sessionID = uuid.New().String()
	}
	createdAt := time.Now().UTC()
	mode := "anonymous"
	if os.Getenv("UPLOG_API_KEY") != "" {
		mode = "authenticated"
	}

	// insert session record
	err = insertSession(sessionID, createdAt, sql.NullTime{}, 0, 0, false, mode,*tag)
	if err != nil {
		fmt.Printf("failed to insert session: %v\n", err)
		os.Exit(1)
	}

	tmpFile := filepath.Join(home, baseDir, tmpDir, sessionID+".log")
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		fmt.Printf("cannot create temp log file: %v\n", err)
		os.Exit(1)
	}
	// Ensure file closed & deleted after finish
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start uploader goroutine that tails tmpFile and uploads new data
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := uploaderLoop(ctx, tmpFile, sessionID,doneWriting)
		if err != nil {
			fmt.Printf("uploader error: %v\n", err)
		}
	}()

	// Run the command
	cmdCtx, cmdCancel := context.WithCancel(context.Background())
	cmdExec := exec.CommandContext(cmdCtx, command, args...)
	// Merge stdout & stderr
	stdoutPipe, err := cmdExec.StdoutPipe()
	if err != nil {
		cancel()
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	stderrPipe, err := cmdExec.StderrPipe()
	if err != nil {
		cancel()
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := cmdExec.Start(); err != nil {
		cancel()
		fmt.Printf("failed to start command: %v\n", err)
		os.Exit(1)
	}

	// reader goroutine: read both streams and write to file
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyStreamsToFile(stdoutPipe, stderrPipe, f, sessionID)
		close(doneWriting)
		// make sure that doneWriting channel is closed, since it would upload the remaining content to the server from .log files
	}()

	// wait for command to finish
	err = cmdExec.Wait()
	// stop uploader after command exits
	cmdCancel()
	// give uploader a moment to flush; then cancel context if still running
	// (uploaderLoop watches ctx)
	time.Sleep(250*time.Millisecond)
	cancel()
	// wait goroutines
	wg.Wait()

	// final mark uploaded
	if err != nil {
		fmt.Printf("command finished with error: %v\n", err)
	} else {
		// mark as uploaded; uploaderLoop should have uploaded already but mark anyway
		if err2 := markSessionUploaded(sessionID); err2 != nil {
			fmt.Printf("warning: failed to mark uploaded: %v\n", err2)
		}
	}

	// delete temp file on completion
	_ = os.Remove(tmpFile)
	fmt.Printf("session finished: %s -> %s\n", sessionID, constructShareURL(sessionID))
}

// listCmd reads sessions from sqlite and prints them
func listCmd(home string) {
	sessions, err := listSessions()
	if err != nil {
		fmt.Printf("failed to list sessions: %v\n", err)
		os.Exit(1)
	}
	if len(sessions) == 0 {
		fmt.Println("no sessions found")
		return
	}
	for _, s := range sessions {
		ts := s.CreatedAt.Format("2006-01-02 15:04:05")
		uploaded := "no"
		if s.IsUploaded {
			uploaded = "yes"
		}
		fmt.Printf("%s | %s | %d bytes | %d lines | uploaded: %s | %s\n",
			ts, s.Tag, s.SizeBytes, s.LineCount, uploaded, constructShareURL(s.SessionID))
	}
}
func deleteCmd(home, sessionID string) {
	if sessionID == "--all" {
		db.Exec(`DELETE FROM sessions`)
		files, _ := filepath.Glob(filepath.Join(home, baseDir, tmpDir, "*.log"))
		for _, f := range files {
			_ = os.Remove(f)
		}
		fmt.Println("deleted all uplog sessions")
		return
	}
	// delete from sqlite
	_, err := db.Exec(`DELETE FROM sessions WHERE session_id = ?`, sessionID)
	if err != nil {
		fmt.Printf("failed to delete session: %v\n", err)
		os.Exit(1)
	}

	// 2. remove temp log file if exists
	tmpPath := filepath.Join(home, baseDir, tmpDir, sessionID+".log")
	_ = os.Remove(tmpPath) // ignore error if not exists

	fmt.Printf("deleted session %s\n", sessionID)
}

// -------------------- File copy & upload --------------------

func copyStreamsToFile(stdout io.ReadCloser, stderr io.ReadCloser, outFile *os.File, sessionID string) {
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
			entry := LogEntry{
				Message:   line,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Level:     detectLevel(line),
			}

			// Serialize to JSONL or a custom format
			b, _ := json.Marshal(entry)
			outFile.Write(b)
			outFile.WriteString("\n")

			incrementSessionStats(sessionID, 1, int64(len(b)))
			// flush to disk for durability
			_ = outFile.Sync()
		}
	}
	go writeLines(stdoutScanner)
	go writeLines(stderrScanner)
	wg.Wait()
}


func uploaderLoop(ctx context.Context, path string, sessionID string, doneWriting <-chan struct{}) error {
	var offset int64 = 0
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		// uploader should stop only when writer has completely finished writing logs
		case <-doneWriting:
			// final drain
			if err := uploadNewChunks(path, &offset, sessionID); err != nil {
				return err
			}
			return nil

		// regular timed uploads
		case <-ticker.C:
			if err := uploadNewChunks(path, &offset, sessionID); err != nil {
				fmt.Printf("upload error (will retry): %v\n", err)
			}

		// context cancel should NOT stop uploader
		// it only stops polling but uploader continues until doneWriting closes
		case <-ctx.Done():
			// do nothing — wait for doneWriting
		}
	}
}


func uploadNewChunks(path string, offset *int64, sessionID string) error {
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

	batch := make([]LogEntry, 0, batchLimit)

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		*offset += int64(len(lineBytes)) + 1 // +1 for newline

		var entry LogEntry
		if err := json.Unmarshal(lineBytes, &entry); err != nil {
			continue
		}

		batch = append(batch, entry)

		if len(batch) >= batchLimit {
			if err := sendBatch(sessionID, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	// Send leftover batch
	if len(batch) > 0 {
		return sendBatch(sessionID, batch)
	}

	return scanner.Err()
}

// sendBatch sends a chunk to the backend
func sendBatch(sessionID string, batch []LogEntry) error {
	body, _ := json.Marshal(batch)
	req, err := http.NewRequest("POST",
		backendUploadEndpoint+"?session_id="+sessionID+"&"+"tag="+*tag,
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

// -------------------- DB helpers --------------------

func insertSession(sessionID string, createdAt time.Time, expiredAt sql.NullTime, lines int64, bytes int64, isUploaded bool, mode string,tag string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
	INSERT OR REPLACE INTO sessions(session_id, created_at, expired_at, line_count, size_bytes, is_uploaded, mode, tag)
	VALUES(?,?,?,?,?,?,?,?)
	`, sessionID, createdAt.Format(time.RFC3339), nullTimeToString(expiredAt), lines, bytes, boolToInt(isUploaded), mode, tag)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func incrementSessionStats(sessionID string, addLines int64, addBytes int64) error {
	_, err := db.Exec(`
	UPDATE sessions SET
		line_count = line_count + ?,
		size_bytes = size_bytes + ?
	WHERE session_id = ?
	`, addLines, addBytes, sessionID)
	return err
}

func markSessionUploaded(sessionID string) error {
	_, err := db.Exec(`UPDATE sessions SET is_uploaded = 1 WHERE session_id = ?`, sessionID)
	return err
}

func listSessions() ([]Session, error) {
	rows, err := db.Query(`SELECT session_id, created_at, expired_at, line_count, size_bytes, is_uploaded, mode, tag FROM sessions ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Session
	for rows.Next() {
		var s Session
		var createdAtStr string
		var expiredAtStr sql.NullString
		var isUploadedInt int
		if err := rows.Scan(&s.SessionID, &createdAtStr, &expiredAtStr, &s.LineCount, &s.SizeBytes, &isUploadedInt, &s.Mode, &s.Tag); err != nil {
			return nil, err
		}
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		if expiredAtStr.Valid {
			t, _ := time.Parse(time.RFC3339, expiredAtStr.String)
			s.ExpiredAt = sql.NullTime{Time: t, Valid: true}
		} else {
			s.ExpiredAt = sql.NullTime{Valid: false}
		}
		s.IsUploaded = isUploadedInt != 0
		res = append(res, s)
	}
	return res, nil
}

// -------------------- Utilities --------------------

func runCmdWithFlags(home string, args []string) {
    fs := flag.NewFlagSet("uplog run", flag.ExitOnError)

    // poll is an int (in milliseconds)
    poll := fs.Int("poll", int(pollIntervalLimit/time.Millisecond), "Polling interval in milliseconds")
	tag = fs.String("tag","","Optional tag name, to tag the session")
    fs.Parse(args)

    if fs.NArg() < 1 {
        fmt.Println("usage: uplog run [--poll N] <command> [args...]")
        os.Exit(1)
    }

    // convert ms → time.Duration
    pollInterval = time.Duration(*poll) * time.Millisecond

    // enforce minimum poll interval (to avoid hammering backend)
    if pollInterval < pollIntervalLimit {
        fmt.Printf("Cannot poll below %v ms. Try --poll >= %v.\n",
            pollIntervalLimit/time.Millisecond,
            pollIntervalLimit/time.Millisecond)
        os.Exit(1)
    }

    command := fs.Arg(0)
    commandArgs := fs.Args()[1:]

    runCmd(home, command, commandArgs)
}


func detectLevel(line string) string {
	up := strings.ToUpper(line)

	switch {
	case strings.Contains(up, "ERROR"):
		return "ERROR"
	case strings.Contains(up, "WARN"):
		return "WARN"
	case strings.Contains(up, "DEBUG"):
		return "DEBUG"
	case strings.Contains(up, "INFO"):
		return "INFO"
	default:
		return "INFO"
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullTimeToString(n sql.NullTime) interface{} {
	if n.Valid {
		return n.Time.Format(time.RFC3339)
	}
	return nil
}

func constructShareURL(sessionID string) string {
	return logsDomain + sessionID
}

// placeholder: ask backend for session id
// Replace with your real API call to create a session. For now return UUID.
func requestSessionIDFromBackend() (string, error) {
	// Example: do POST to backend create session, pass auth header if present, parse returned id
	// For now we do a best-effort call to /create-session; if it fails, return error and caller will fallback to local UUID.
	// Try a quick request to backend (non-fatal)
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("POST", backendSessionCreateEndpoint, nil)
	if err != nil {
		return "", err
	}
	if key := os.Getenv("UPLOG_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// read body as id (assume plain text or JSON id)
		json.NewDecoder(resp.Body).Decode(&respJSON)
		sessionID := respJSON.SessionID
		if sessionID != "" {
			return sessionID, nil
		}
	}
	// fallback: return error so caller can use uuid
	return "", fmt.Errorf("backend returned status %d", resp.StatusCode)
}
