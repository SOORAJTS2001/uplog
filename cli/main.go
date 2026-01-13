package main

import (
	"cli/constants"
	"cli/db"
	"cli/executor"
	"cli/setup"
	"cli/tail"
	"cli/utils"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"time"
	"io"
)

var reservedCommands = []string{"list","purge","upload","delete","-h","-help","--h","--help"}

func runArgParser() (int, int, string, []string) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard) // silence default help/errors

	poll := fs.Int("poll", constants.PollIntervalLimit, "Polling interval")
	batch := fs.Int("batch", constants.BatchLimit, "Batch size")
	tag := fs.String("tag", "", "Session tag")

	_ = fs.Parse(os.Args[1:])

	return *poll, *batch, *tag, fs.Args()
}

func reservedArgParser(args []string){
	switch args[1] {
	case "list":utils.ListLogs(setup.DB_con)
	case "purge":utils.PurgeLogs(setup.DB_con)
	case "upload":utils.UploadLog(args[1:])
	case "delete":utils.DeleteLog(args[2],setup.DB_con)
	default: utils.Help()
	}
}


func main() {
	args := os.Args
	command, commandArgs := args[1], args[2:]
	if slices.Contains(reservedCommands,command){
		setup.Setup(false)
		reservedArgParser(args)
		return
	}
	setup.Setup(true)
	pollInterval, batchSize, tag,parsed_args := runArgParser()
	command,commandArgs = parsed_args[0],parsed_args[1:]
	fmt.Println(pollInterval,batchSize,tag)
	db.UpsertSessionById(setup.DB_con,utils.NewSession(setup.SessionId,tag,setup.UserId))
	wg := &sync.WaitGroup{}
	tmpFile := filepath.Join(constants.TmpDir, setup.SessionId+".log")
	os.OpenFile(tmpFile,os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	wg.Add(1)
	go tail.Tail(wg, (time.Duration(pollInterval) * time.Millisecond), batchSize, tag, setup.SessionId, tmpFile)
	// This delay is important so that the tail can get the start and end position
	time.Sleep(50 * time.Millisecond)

	// --- Start executor in foreground ---
	executorArgs := append([]string{command}, commandArgs...)
	fmt.Println(executorArgs)
	wg.Add(1)
	go executor.Executor(wg, executorArgs,setup.SessionId, tmpFile)

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
	os.Remove(tmpFile)
}
