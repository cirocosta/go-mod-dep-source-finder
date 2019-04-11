package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func failWithHelp(messageFormat string, args ...interface{}) {
	const usageFormat = `Usage: %s <line>\n`

	fmt.Fprintf(os.Stderr, messageFormat, args...)
	fmt.Fprintf(os.Stderr, usageFormat, os.Args[0])

	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		failWithHelp("error: not enough arguments.\n")
	}

	var reader io.Reader

	line := os.Args[1]
	if line == "-" {
		reader = os.Stdin
	} else {
		reader = bytes.NewBufferString(line)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("line: %s\n", scanner.Text())
	}
}
