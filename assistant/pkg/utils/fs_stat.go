package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type StatResult struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	IsDir      bool   `json:"is_dir"`
	Mode       string `json:"mode"`
	ModTime    string `json:"mod_time"`
	IsSymlink  bool   `json:"is_symlink"`
	LinkTarget string `json:"link_target,omitempty"`
}

func Stat(path string) (*StatResult, error) {
	return StatInDir(".", path)
}

func StatInDir(baseDir, relPath string) (*StatResult, error) {
	clean, err := cleanRelFSPath(relPath)
	if err != nil {
		return nil, err
	}
	full := filepath.Join(baseDir, clean)

	info, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, clean)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, clean)
		}
		return nil, fmt.Errorf("stat failed: %w", err)
	}

	res := &StatResult{
		Path:    full,
		Name:    info.Name(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		Mode:    info.Mode().String(),
		ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
	}

	if info.Mode()&os.ModeSymlink != 0 {
		res.IsSymlink = true
		if target, err := os.Readlink(full); err == nil {
			res.LinkTarget = target
		}
		// For symlinks, keep Lstat() metadata; do not follow.
	}

	return res, nil
}
