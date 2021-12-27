package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dottedmag/xd"
)

func main() {
	var buf [4096]byte
	var offset int
	for {
		n, err := os.Stdin.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read: %v\n", err)
			os.Exit(1)
		}

		xd.Print(buf[:n], offset)
		offset += n
	}
}
