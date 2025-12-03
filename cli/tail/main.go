package tail

import (
	"bufio"
	"cli/api"
	"cli/models"
	"cli/tail/utils"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strings"
	"time"
	"sync"
)

var batchedLogs []models.LogEntry
var backendDisabled bool = false

func Tail(wg *sync.WaitGroup,pollInterval time.Duration,batchSize int,tag string, sessionId string,filePath string) {
	defer wg.Done()
	fmt.Println("Poll:", pollInterval, "Batch:", batchSize)

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
				if backendDisabled{
					return
				}
				if line == "EOF" {
					api.BatchUpload(batchedLogs, sessionId, tag, 3, &backendDisabled)
					os.Exit(0)
				}
				if len(batchedLogs) == batchSize {
					api.BatchUpload(batchedLogs, sessionId, tag, 3, &backendDisabled)
					batchedLogs = batchedLogs[:0]
				}
				start := fmt.Sprintf("<%s>",sessionId)
				end := fmt.Sprintf("</%s>",sessionId)

				if strings.HasPrefix(line, start) && strings.HasSuffix(line, end) {
					content := line[len(start) : len(line)-len(end)]
					entry := models.LogEntry{
						Message:   content,
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
