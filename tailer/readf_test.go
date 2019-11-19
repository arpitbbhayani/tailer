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

type ReadFTesting struct {
	t        *testing.T
	filepath string
	fp       *os.File
}

func (t *ReadFTesting) Init() {
	t.filepath = "./temp.txt"
	createFile(t.filepath)
	fp, err := os.OpenFile(t.filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		panic(err)
	}
	t.fp = fp
}

func (t *ReadFTesting) End() {
	deleteFile(t.filepath)
	t.fp.Close()
}

func (t *ReadFTesting) Test(actualLinesData []string, readFConfig ReadFConfig) {
	actualLines := make(chan string)
	go func() {
		for _, actualLine := range actualLinesData {
			actualLines <- actualLine
			time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		}
		close(actualLines)
	}()

	go func() {
		WriteF(t.fp, actualLines)
	}()

	observedLines, err := ReadF(t.filepath, readFConfig)
	if err != nil {
		panic(err)
	}
	verifyReadLines(observedLines, actualLinesData, t.t)
}

func TestFixedWidthFullContent(t *testing.T) {
	rt := &ReadFTesting{}
	rt.Init()
	defer rt.End()

	totalActualLines := 100
	actualLinesData := make([]string, 0, totalActualLines)
	for i := 0; i < totalActualLines; i++ {
		actualLinesData = append(actualLinesData, randomString(10))
	}

	rt.Test(actualLinesData, ReadFConfig{
		SizePerMessageInBytes: 10,
		MaxMessagesInBuffer:   10,
	})
}

func TestFixedWidthHalfContent(t *testing.T) {
	rt := &ReadFTesting{}
	rt.Init()
	defer rt.End()

	totalActualLines := 100
	actualLinesData := make([]string, 0, totalActualLines)
	for i := 0; i < totalActualLines; i++ {
		actualLinesData = append(actualLinesData, randomString(5))
	}

	rt.Test(actualLinesData, ReadFConfig{
		SizePerMessageInBytes: 10,
		MaxMessagesInBuffer:   10,
	})
}

func TestFixedWidthRandomContent(t *testing.T) {
	rt := &ReadFTesting{}
	rt.Init()
	defer rt.End()

	totalActualLines := 100
	actualLinesData := make([]string, 0, totalActualLines)
	for i := 0; i < totalActualLines; i++ {
		actualLinesData = append(actualLinesData, randomString(rand.Intn(10)))
	}

	rt.Test(actualLinesData, ReadFConfig{
		SizePerMessageInBytes: 10,
		MaxMessagesInBuffer:   10,
	})
}
