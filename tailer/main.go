package tailer

import (
	"os"
	"strings"
	"time"
)

func Tailf(filePath string) (<-chan string, error) {
	buffer := make([]byte, 4096)
	ch := make(chan string)

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			bytesRead, _ := fp.Read(buffer)
			if bytesRead > 0 {
				for _, line := range strings.Split(string(buffer[:bytesRead]), "\n") {
					ch <- line
				}
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return ch, nil
}
