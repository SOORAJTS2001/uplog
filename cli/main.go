package main

import (
	"cli/api"
	"cli/constants"
	"cli/executor"
	"cli/setup"
	"cli/tail"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func argParser() (int, int, string) {
	poll := flag.Int("poll", constants.PollIntervalLimit, "Default polling time in milliseconds")
	batchSize := flag.Int("batch", constants.BatchLimit, "Default polling batch")
	tag := flag.String("tag", "", "Tag for this session")
	return *poll, *batchSize, *tag
}

func main() {
	setup.Setup()
	pollInterval, batchSize, tag := argParser()
	sessionId, err := api.SetupSession()
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	tmpFile := filepath.Join(constants.HomeDir, constants.BaseDir, constants.TmpDir, sessionId+".log")
	os.OpenFile(tmpFile,os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	args := os.Args
	command, commandArgs := args[1], args[2:]
	wg.Add(1)
	go tail.Tail(wg, (time.Duration(pollInterval) * time.Millisecond), batchSize, tag, sessionId, tmpFile)
	// This delay is important so that the tail can get the start and end position
	time.Sleep(50 * time.Millisecond)

	// --- Start executor in foreground ---
	executorArgs := append([]string{command}, commandArgs...)
	fmt.Println(executorArgs)
	wg.Add(1)
	go executor.Executor(wg, executorArgs,sessionId, tmpFile)

	// --- SIGNAL HANDLING ---
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nStopping both processes...")
		os.Exit(1)
	}()

	// Wait for executor to finish
	fmt.Println("Executor finished.")

	// Wait for tailer to finish
	wg.Wait()
	fmt.Println("Tailer finished.")
	fmt.Println("Wrapper finished.")
	wg.Done()
	os.Remove(tmpFile)
}
