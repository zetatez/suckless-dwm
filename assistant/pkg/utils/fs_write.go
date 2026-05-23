package utils

import (
	"fmt"
	"os"
)

type WriteMode string

const (
	WriteModeOverwrite WriteMode = "overwrite"
	WriteModeAppend    WriteMode = "append"
	WriteModeCreate    WriteMode = "create"
)

func Write(path string, content string, mode WriteMode) error {
	switch mode {
	case WriteModeOverwrite:
		return os.WriteFile(path, []byte(content), 0644)
	case WriteModeAppend:
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open file failed: %w", err)
		}
		defer f.Close()
		if _, err := f.WriteString(content); err != nil {
			return fmt.Errorf("append failed: %w", err)
		}
		return nil
	case WriteModeCreate:
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists: %s", path)
		}
		return os.WriteFile(path, []byte(content), 0644)
	default:
		return fmt.Errorf("unknown write mode: %s", mode)
	}
}
