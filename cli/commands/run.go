package commands

import (
	"fmt"
	"path/filepath"
	"time"
	"os"
	"sync"
	"context"
	"os/exec"
	"database/sql"
	"github.com/google/uuid"
	"uplog/utils"
	"uplog/config"
	"uplog/workers"
	"uplog/db"
	"uplog/models"
)

func RunCmd(sql_object *sql.DB,home, command string, args []string,respJSON models.SessionCreateResponse,pollInterval time.Duration,tag *string) {
	// 1. Get session ID from backend (or local generation if backend not available)
	sessionID, err := utils.RequestSessionIDFromBackend(respJSON)
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
	err = db.InsertSession(sql_object,sessionID, createdAt, sql.NullTime{}, 0, 0, false, mode,*tag)
	if err != nil {
		fmt.Printf("failed to insert session: %v\n", err)
		os.Exit(1)
	}

	tmpFile := filepath.Join(home, config.BaseDir, config.TmpDir, sessionID+".log")
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
		err := worker.UploaderLoop(ctx, tmpFile, sessionID,pollInterval,tag,doneWriting)
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
		worker.CopyStreamsToFile(sql_object,stdoutPipe, stderrPipe, f, sessionID)
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
		if err2 := db.MarkSessionUploaded(sql_object,sessionID); err2 != nil {
			fmt.Printf("warning: failed to mark uploaded: %v\n", err2)
		}
	}

	// delete temp file on completion
	_ = os.Remove(tmpFile)
	fmt.Printf("session finished: %s -> %s\n", sessionID, utils.ConstructShareURL(sessionID))
}
