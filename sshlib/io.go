package sshlib

import (
	"bufio"
	"io"
	"regexp"
)

type WriteCloser struct {
	io.WriteCloser
}

func (w *WriteCloser) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w *WriteCloser) WriteLine(s string) (n int, err error) {
	return w.Write([]byte(s + "\r"))
}

type BufferedReader struct {
	*bufio.Reader
}

// NewBufferedReader returns a BufferedReader that writes to wList what it reads from r.
func NewBufferedReader(r io.Reader, wList ...io.Writer) *BufferedReader {
	return &BufferedReader{Reader: bufio.NewReader(TeeReader(r, wList...))}
}

// ReadUntil reads from r one byte at a time until the bytes read match the given pattern
func (br BufferedReader) ReadUntil(pattern string) (result []byte, err error) {
	for {
		b, _ := br.ReadByte()
		result = append(result, b)
		matched, err := regexp.Match(pattern, result)
		if matched || err != nil {
			return result, err
		}
	}
}

// TeeReader returns an io.Reader that writes to wList what it reads from r.
func TeeReader(r io.Reader, wList ...io.Writer) io.Reader {
	mw := io.MultiWriter(wList...)
	return io.TeeReader(r, mw)
}
