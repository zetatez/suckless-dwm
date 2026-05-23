package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// cleanRelFSPath validates a relative path intended to be resolved under a base dir.
// It rejects empty paths, absolute paths, '.', and any path that would escape the base.
func cleanRelFSPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", fmt.Errorf("%w: empty path", ErrInvalidPath)
	}
	if filepath.IsAbs(p) {
		return "", fmt.Errorf("%w: absolute path not allowed: %s", ErrInvalidPath, p)
	}
	clean := filepath.Clean(p)
	if clean == "." {
		return "", fmt.Errorf("%w: invalid path: %s", ErrInvalidPath, p)
	}
	if strings.HasPrefix(clean, ".."+string(filepath.Separator)) || clean == ".." {
		return "", fmt.Errorf("%w: path escapes base dir: %s", ErrInvalidPath, p)
	}
	return clean, nil
}
