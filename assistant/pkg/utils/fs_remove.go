package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RemoveResult struct {
	Path    string `json:"path"`
	Removed bool   `json:"removed"`
}

type RemoveCmd struct {
	baseDir   string
	paths     []string
	recursive bool
	force     bool
}

func NewRemove(paths ...string) *RemoveCmd {
	return &RemoveCmd{baseDir: ".", paths: paths}
}

func (r *RemoveCmd) InDir(baseDir string) *RemoveCmd {
	if strings.TrimSpace(baseDir) != "" {
		r.baseDir = baseDir
	}
	return r
}

func (r *RemoveCmd) WithRecursive() *RemoveCmd {
	r.recursive = true
	return r
}

func (r *RemoveCmd) WithForce() *RemoveCmd {
	r.force = true
	return r
}

func (r *RemoveCmd) Exec() ([]RemoveResult, error) {
	if len(r.paths) == 0 {
		return nil, fmt.Errorf("%w: no paths", ErrInvalidPath)
	}

	base := r.baseDir
	if strings.TrimSpace(base) == "" {
		base = "."
	}

	results := make([]RemoveResult, 0, len(r.paths))
	for _, p := range r.paths {
		clean, err := cleanRelFSPath(p)
		if err != nil {
			return nil, err
		}
		full := filepath.Join(base, clean)

		info, err := os.Lstat(full)
		if err != nil {
			if os.IsNotExist(err) {
				if r.force {
					results = append(results, RemoveResult{Path: full, Removed: false})
					continue
				}
				return nil, fmt.Errorf("%w: %s", ErrFileNotFound, clean)
			}
			if os.IsPermission(err) {
				return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, clean)
			}
			return nil, fmt.Errorf("stat failed: %w", err)
		}

		if info.IsDir() && !r.recursive {
			return nil, fmt.Errorf("refusing to remove directory without recursive: %s", clean)
		}

		if info.IsDir() {
			if err := os.RemoveAll(full); err != nil {
				return nil, fmt.Errorf("remove dir failed: %w", err)
			}
		} else {
			if err := os.Remove(full); err != nil {
				return nil, fmt.Errorf("remove file failed: %w", err)
			}
		}
		results = append(results, RemoveResult{Path: full, Removed: true})
	}

	return results, nil
}
