package utils

import (
	"bufio"
	"io"
	"strings"
)

// LineScanner is a minimal line iterator without bufio.Scanner token limits.
// It keeps allocations low by reusing its internal buffer.
type LineScanner struct {
	r    *bufio.Reader
	line string
	err  error
}

func NewLineScanner(r io.Reader) *LineScanner {
	return &LineScanner{r: bufio.NewReaderSize(r, 64*1024)}
}

func (s *LineScanner) Scan() bool {
	if s.err != nil {
		return false
	}

	line, err := s.r.ReadString('\n')
	if err == io.EOF {
		if len(line) == 0 {
			s.err = io.EOF
			return false
		}
		// last line without trailing newline
		s.line = strings.TrimRight(line, "\r\n")
		return true
	}
	if err != nil {
		s.err = err
		return false
	}
	// Strip trailing newline.
	s.line = strings.TrimRight(line, "\r\n")
	return true
}

func (s *LineScanner) Text() string { return s.line }

func (s *LineScanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
