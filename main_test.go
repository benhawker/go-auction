package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_e2e_Success(t *testing.T) {
	output := captureOutput(func() {
		parseFile(inputFile)
	})

	expectedOutputLines := []string{
		"20|toaster_1|8|SOLD|12.50|3|20.00|7.50",
		"20|tv_1| |UNSOLD|0.00|2|200.00|150.00",
	}

	for _, eol := range expectedOutputLines {
		if strings.Contains(output, eol) != true {
			t.Errorf("Expected %v to contain %v.", output, eol)
		}
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
