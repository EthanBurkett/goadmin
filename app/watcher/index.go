package watcher

import (
	"fmt"
	"os"
	"time"
)

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error accessing file:", err)
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error accessing file:", err)
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			fmt.Println("File has been modified")
			initialStat = stat
		}

		time.Sleep(500 * time.Millisecond)
	}
}
