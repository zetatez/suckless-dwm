package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
}

func ListDir(path string) ([]FileInfo, error) {
	if strings.TrimSpace(path) == "" {
		path = "."
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir failed: %w", err)
	}

	results := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("stat entry failed (%s): %w", entry.Name(), err)
		}
		results = append(results, FileInfo{
			Path:    filepath.Join(path, entry.Name()),
			Name:    entry.Name(),
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	return results, nil
}
