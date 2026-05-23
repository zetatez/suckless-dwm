package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReadResult struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Lines    int    `json:"lines"`
	Encoding string `json:"encoding"`
}

func (f *ReadResult) Bytes() []byte {
	return []byte(f.Content)
}

func (f *ReadResult) Reader() io.Reader {
	return bytes.NewReader([]byte(f.Content))
}

func Read(path string) (*ReadResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	return &ReadResult{
		Path:     path,
		Content:  string(content),
		Lines:    strings.Count(string(content), "\n") + 1,
		Encoding: detectEncoding(content),
	}, nil
}

func ReadLines(path string, start, end int) (*ReadResult, error) {
	if start < 1 || end < start {
		return nil, fmt.Errorf("%w: start=%d, end=%d", ErrInvalidRange, start, end)
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	lineNum := 0
	var lines []string
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if lineNum >= start && lineNum <= end {
			lines = append(lines, fmt.Sprintf("%6d | %s", lineNum, line))
		}
		if lineNum >= end {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file failed: %w", err)
	}

	return &ReadResult{
		Path:     path,
		Content:  strings.Join(lines, "\n"),
		Lines:    lineNum,
		Encoding: "utf-8",
	}, nil
}

func ReadChunk(path string, offset, size int) (*ReadResult, error) {
	if offset < 0 || size <= 0 {
		return nil, fmt.Errorf("invalid offset or size")
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, path)
		}
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	if _, err := file.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek file failed: %w", err)
	}

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	return &ReadResult{
		Path:     path,
		Content:  string(buf[:n]),
		Encoding: "binary",
	}, nil
}
