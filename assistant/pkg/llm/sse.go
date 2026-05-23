package llm

import (
	"bufio"
	"context"
	"io"
	"strings"
)

// ReadSSE reads Server-Sent Events from r and calls onData for each event's
// assembled data field.
//
// It supports multi-line data fields and ignores other SSE fields.
func ReadSSE(ctx context.Context, r io.Reader, onData func(data string) error) error {
	scanner := bufio.NewScanner(r)
	// Default scanner token limit is 64K; streaming chunks can be larger.
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	dataLines := make([]string, 0, 8)
	dispatch := func() error {
		if len(dataLines) == 0 {
			return nil
		}
		data := strings.Join(dataLines, "\n")
		dataLines = dataLines[:0]
		if strings.TrimSpace(data) == "" {
			return nil
		}
		return onData(data)
	}

	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}

		line := scanner.Text()
		if line == "" {
			if err := dispatch(); err != nil {
				return err
			}
			continue
		}
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "data:") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			dataLines = append(dataLines, v)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return dispatch()
}
