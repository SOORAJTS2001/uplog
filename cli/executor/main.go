package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"github.com/creack/pty"
)

func Executor(wg *sync.WaitGroup,args []string,sessionId string,filename string) {
	defer wg.Done()
	// Replace with any program you want
	cmd := exec.Command(args[0], args[1:]...)

	// Start command inside a PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer ptmx.Close()
	// 2. Open *same file* for appending
	logFile, err := os.OpenFile(filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Read PTY output live
	scanner := bufio.NewScanner(ptmx)
	start := fmt.Sprintf("<%s>",sessionId)
	end := fmt.Sprintf("</%s>",sessionId)
	for scanner.Scan() {
		line := scanner.Text()
		// this is the live logs from the programme, don't comment it
		fmt.Println(line)

		line = start + line + end + "\n"

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
