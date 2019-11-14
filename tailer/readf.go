package tailer

import (
	"bytes"
	"io"
	"os"
	"time"
)

func readMessages(fp *os.File, buf []byte) (int, int, error) {
	bytesRead, err := fp.Read(buf)
	if err != nil {
		return 0, 0, err
	}
	index := bytes.LastIndex(buf[:bytesRead], []byte{'\n'})
	if index == -1 {
		return bytesRead, 0, nil
	}
	return index, (bytesRead - 1) - index, nil
}

// ReadF Reads messages from a ReadOnly file separated by \n
// Reads at max 100 messages and keeps it available in buffer
// to fast fanout.
func ReadF(filePath string) (<-chan string, error) {
	sizePerMessageInBytes, maxMessages := 10+1, 5

	ch := make(chan string, maxMessages)
	buffer := make([]byte, sizePerMessageInBytes*maxMessages)

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}

	go func() {
		defer fp.Close()
		for {
			bytesRead, seekBack, err := readMessages(fp, buffer)
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				fp.Close()
				panic(err)
			}
			fp.Seek(-int64(seekBack), 1)
			for _, message := range bytes.Split(buffer[:bytesRead], []byte{'\n'}) {
				ch <- string(message)
			}
		}
	}()

	return ch, nil
}
