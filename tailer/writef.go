package tailer

import (
	"os"
)

// WriteF writes messages to an appendonly file. Messages are separated by \n
func WriteF(filePath string, ch <-chan string) error {
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	for line := range ch {
		_, err := fp.WriteString(line + "\n")
		if err != nil {
			return err
		}
		fp.Sync()
	}

	return nil
}
