package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	for _, fileName := range os.Args[1:] {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if _, err := io.Copy(os.Stdout, file); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		file.Close()
	}
}
