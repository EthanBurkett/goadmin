package watcher

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ethanburkett/goadmin/app/config"
)

type FileChangeEvent struct {
	FilePath string
	NewLine  string
	ModTime  time.Time
}

func WatchGamesMp(config *config.Config) <-chan FileChangeEvent {
	changesChan := make(chan FileChangeEvent)

	go func() {
		defer close(changesChan)

		initialStat, err := os.Stat(config.GamesMpPath)
		if err != nil {
			fmt.Println("Error accessing file:", err)
			return
		}

		var lastSize int64 = initialStat.Size()
		var lastModTime time.Time = initialStat.ModTime()

		for {
			stat, err := os.Stat(config.GamesMpPath)
			if err != nil {
				fmt.Println("Error accessing file:", err)
				return
			}

			// Check if size changed OR modification time changed
			if stat.Size() != lastSize || stat.ModTime() != lastModTime {

				newLine, err := readLastLine(config.GamesMpPath)
				if err == nil && newLine != "" {

					changesChan <- FileChangeEvent{
						FilePath: config.GamesMpPath,
						NewLine:  newLine,
						ModTime:  stat.ModTime(),
					}
				} else if err != nil {
					fmt.Printf("Error reading last line: %v\n", err)
				}
				lastSize = stat.Size()
				lastModTime = stat.ModTime()
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return changesChan
}

func readLastLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}
