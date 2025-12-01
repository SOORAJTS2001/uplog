package main

import (
	"bufio"
	"cli/setup"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func argParser() (int, int,string) {
	poll := flag.Int("poll", 200, "Default polling time in milliseconds")
	batchSize := flag.Int("batch", 200, "Default polling batch")
	tag:=flag.String("tag","","Tag for this session")
	return *poll, *batchSize,*tag
}

// Streams stdout+stderr of a command live to terminal
func streamOutput(prefix string, r io.ReadCloser) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}
func main() {
	setup.Setup()
	pollInterval,batchSize,tag:=argParser()
	args := os.Args
	command, commandArgs := args[1], args[2:]

	// --- Start tailer in background ---
	tailCmd := exec.Command("go", "run", "tail/main.go","--poll",fmt.Sprint(pollInterval),"--batch",fmt.Sprint(batchSize),"--tag",tag)

	tailStdout, err := tailCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	tailCmd.Stderr = tailCmd.Stdout

	if err := tailCmd.Start(); err != nil {
		fmt.Println("Failed to start tailer:", err)
		return
	}

	fmt.Println("Tailer started:", tailCmd.Process.Pid)

	// Tailer output goroutine
	go streamOutput("TAIL", tailStdout)
	time.Sleep(50*time.Millisecond)

	// --- Start executor in foreground ---
	executorArgs := append([]string{"run", "executor/main.go", command}, commandArgs...)
	fmt.Println(executorArgs)
	execCmd := exec.Command("go", executorArgs...)

	execStdout, err := execCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	execCmd.Stderr = execCmd.Stdout

	if err := execCmd.Start(); err != nil {
		fmt.Println("Failed to start executor:", err)
		return
	}

	fmt.Println("Executor started:", execCmd.Process.Pid)

	// Executor output goroutine
	go streamOutput("EXEC", execStdout)

	// --- SIGNAL HANDLING ---
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nStopping both processes...")
		tailCmd.Process.Kill()
		execCmd.Process.Kill()
	}()

	// Wait for executor to finish
	execCmd.Wait()
	fmt.Println("Executor finished.")

	// Wait for tailer to finish
	tailCmd.Wait()
	fmt.Println("Tailer finished.")

	fmt.Println("Wrapper finished.")
}
