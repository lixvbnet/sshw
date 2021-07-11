package sshlib

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestReadUntil(t *testing.T) {
	r := strings.NewReader("This is some text\n")
	br := NewBufferedReader(r)

	result, err := br.ReadUntil("te")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != "This is some te" {
		t.Fail()
	}
}

func TestTeeReader(t *testing.T) {
	text := "This is some text\n"
	r := strings.NewReader(text)

	fileName := "testTeeReader.output"
	file, _ := os.Create(fileName)
	defer file.Close(); os.Remove(fileName)

	//tr := TeeReader(r, os.Stdout, file)
	tr := TeeReader(r, file)
	buf := make([]byte, 4)
	bytesRead := 0
	for {
		n, err := tr.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		bytesRead += n
	}

	if bytesRead != len(text) {
		t.Fatalf("%d out of %d read\n", bytesRead, len(text))
	}
}
