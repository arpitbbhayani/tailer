package tailer

import (
	"os"
	"strings"
	"time"
)

func TailF(filePath string) (<-chan string, error) {
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
				contentBuffer := buffer[:bytesRead]
				if contentBuffer[bytesRead-1] == byte('\n') {
					contentBuffer = contentBuffer[:bytesRead-1]
				}
				for _, line := range strings.Split(string(contentBuffer), "\n") {
					ch <- line
				}
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return ch, nil
}
