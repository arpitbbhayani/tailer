package tailer

import (
	"os"
)

// WriteF writes messages to an appendonly file. Messages are separated by \n
func WriteF(fp *os.File, ch <-chan string) error {
	for line := range ch {
		_, err := fp.WriteString(line + "\n")
		if err != nil {
			return err
		}
		fp.Sync()
	}
	return nil
}
