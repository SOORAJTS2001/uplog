package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	// Replace with any program you want
	args := os.Args
	fmt.Println(args)
	fmt.Println(args[1], args[2:])
	cmd := exec.Command(args[1], args[2:]...)

	// Start command inside a PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer ptmx.Close()

	// 1. Truncate (recreate) the file
	f, err := os.OpenFile("sample.log",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	f.Close() // important — close the truncate handle

	// 2. Open *same file* for appending
	logFile, err := os.OpenFile("sample.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Read PTY output live
	scanner := bufio.NewScanner(ptmx)
	for scanner.Scan() {
		line := scanner.Text()
		line = "<uplog>" + line + "</uplog>\n"

		fmt.Println(line)
		// Write live to file
		lineBytes := []byte(line)
		_, err = logFile.Write(lineBytes)
		if err != nil {
			panic(err)
		}

		logFile.Sync()
	}
	logFile.WriteString("EOF")
	logFile.Sync()
	cmd.Wait()
	fmt.Println("process exited.")
}
