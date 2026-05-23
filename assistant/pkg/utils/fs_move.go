package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type MoveResult struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Move struct {
	baseDir   string
	from      string
	to        string
	overwrite bool
}

func NewMove(from, to string) *Move {
	return &Move{baseDir: ".", from: from, to: to}
}

func (m *Move) InDir(baseDir string) *Move {
	if strings.TrimSpace(baseDir) != "" {
		m.baseDir = baseDir
	}
	return m
}

func (m *Move) WithOverwrite() *Move {
	m.overwrite = true
	return m
}

func (m *Move) Exec() (*MoveResult, error) {
	fromRel, err := cleanRelFSPath(m.from)
	if err != nil {
		return nil, err
	}
	toRel, err := cleanRelFSPath(m.to)
	if err != nil {
		return nil, err
	}
	fromFull := filepath.Join(m.baseDir, fromRel)
	toFull := filepath.Join(m.baseDir, toRel)

	if !Exists(fromFull) {
		return nil, fmt.Errorf("%w: %s", ErrFileNotFound, fromRel)
	}
	if Exists(toFull) && !m.overwrite {
		return nil, fmt.Errorf("destination exists: %s", toRel)
	}
	if err := os.MkdirAll(filepath.Dir(toFull), 0755); err != nil {
		return nil, fmt.Errorf("mkdir failed: %w", err)
	}
	if Exists(toFull) && m.overwrite {
		if err := os.RemoveAll(toFull); err != nil {
			return nil, fmt.Errorf("remove destination failed: %w", err)
		}
	}

	if err := os.Rename(fromFull, toFull); err == nil {
		return &MoveResult{From: fromFull, To: toFull}, nil
	} else {
		var linkErr *os.LinkError
		if errors.As(err, &linkErr) {
			if errno, ok := linkErr.Err.(syscall.Errno); ok && errno == syscall.EXDEV {
				if err := moveByCopy(fromFull, toFull); err != nil {
					return nil, err
				}
				return &MoveResult{From: fromFull, To: toFull}, nil
			}
		}
		return nil, fmt.Errorf("move failed: %w", err)
	}
}

func moveByCopy(fromFull, toFull string) error {
	info, err := os.Lstat(fromFull)
	if err != nil {
		return fmt.Errorf("stat source failed: %w", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(fromFull)
		if err != nil {
			return fmt.Errorf("readlink failed: %w", err)
		}
		if err := os.Symlink(target, toFull); err != nil {
			return fmt.Errorf("symlink failed: %w", err)
		}
		if err := os.Remove(fromFull); err != nil {
			return fmt.Errorf("remove source symlink failed: %w", err)
		}
		return nil
	}

	if info.IsDir() {
		if err := copyDir(fromFull, toFull); err != nil {
			return err
		}
		if err := os.RemoveAll(fromFull); err != nil {
			return fmt.Errorf("remove source dir failed: %w", err)
		}
		return nil
	}

	if err := copyFile(fromFull, toFull, info.Mode()); err != nil {
		return err
	}
	_ = os.Chtimes(toFull, info.ModTime(), info.ModTime())
	if err := os.Remove(fromFull); err != nil {
		return fmt.Errorf("remove source file failed: %w", err)
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source failed: %w", err)
	}
	defer in.Close()

	dir := filepath.Dir(dst)
	base := filepath.Base(dst)
	tmp, err := os.CreateTemp(dir, base+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if err := tmp.Chmod(mode.Perm()); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temp file failed: %w", err)
	}
	if _, err := io.Copy(tmp, in); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("copy failed: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temp file failed: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file failed: %w", err)
	}

	if err := os.Rename(tmpName, dst); err != nil {
		return fmt.Errorf("rename temp file failed: %w", err)
	}
	return nil
}

func copyDir(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("stat source dir failed: %w", err)
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return fmt.Errorf("mkdir dest dir failed: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read source dir failed: %w", err)
	}
	for _, e := range entries {
		srcP := filepath.Join(src, e.Name())
		dstP := filepath.Join(dst, e.Name())

		ei, err := os.Lstat(srcP)
		if err != nil {
			return fmt.Errorf("stat entry failed: %w", err)
		}

		if ei.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(srcP)
			if err != nil {
				return fmt.Errorf("readlink failed: %w", err)
			}
			if err := os.Symlink(target, dstP); err != nil {
				return fmt.Errorf("symlink failed: %w", err)
			}
			continue
		}

		if ei.IsDir() {
			if err := copyDir(srcP, dstP); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(srcP, dstP, ei.Mode()); err != nil {
			return err
		}
		_ = os.Chtimes(dstP, ei.ModTime(), ei.ModTime())
	}
	return nil
}
