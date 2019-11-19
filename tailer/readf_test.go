package tailer

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return randomStringWithCharset(length, charset)
}

func createFile(filepath string) {
	_, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
}

func deleteFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		panic(err)
	}
}

func verifyReadLines(ch <-chan string, actualLines []string, t *testing.T) {
	linesVerified := 0
	for _, actualLine := range actualLines {
		observedLine := <-ch
		if observedLine != actualLine {
			t.Error("observed and actual lines do not match", observedLine, actualLine)
		} else {
			linesVerified++
		}
	}
	if linesVerified != len(actualLines) {
		t.Error("number of lines verified is", linesVerified, "but total lines to be verified were", len(actualLines))
	}
}

func TestFixedWidthContent(t *testing.T) {
	filepath := "./temp.txt"
	createFile(filepath)
	defer deleteFile(filepath)

	totalActualLines := 100
	actualLines := make(chan string)
	actualLinesData := make([]string, 0, totalActualLines)

	for i := 0; i < totalActualLines; i++ {
		actualLinesData = append(actualLinesData, randomString(10))
	}

	go func() {
		for _, actualLine := range actualLinesData {
			actualLines <- actualLine
			time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		}
		close(actualLines)
	}()

	go func() {
		WriteF(filepath, actualLines)
	}()

	observedLines, err := ReadF(filepath, ReadFConfig{
		SizePerMessageInBytes: 10,
		MaxMessagesInBuffer:   10,
	})
	if err != nil {
		panic(err)
	}
	verifyReadLines(observedLines, actualLinesData, t)
}
