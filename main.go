package main

import (
	"fmt"
	"os"

	"github.com/arpitbbhayani/tailer/tailer"
)

func main() {
	filePath := os.Args[1]

	ch, err := tailer.Tailf(filePath)
	if err != nil {
		panic(err)
	}
	for line := range ch {
		fmt.Println(line)
	}
}
