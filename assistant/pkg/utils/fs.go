package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func IsDirExists(path string) bool {
	return IsDir(path)
}

func IsFileExists(path string) bool {
	return IsFile(path)
}

func GetAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func Mkdir(path string, parents bool) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("mkdir failed: empty path")
	}
	if parents {
		return os.MkdirAll(path, 0755)
	}
	return os.Mkdir(path, 0755)
}

func Copy(src, dst string) error {
	return CopyFile(src, dst)
}

func CopyFile(src, dst string) error {
	if strings.TrimSpace(src) == "" || strings.TrimSpace(dst) == "" {
		return fmt.Errorf("copy failed: empty src or dst")
	}
	if src == dst {
		return nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source failed: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source failed: %w", err)
	}
	if srcInfo.IsDir() {
		return fmt.Errorf("copy failed: source is a directory: %s", src)
	}

	if dir := filepath.Dir(dst); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create destination dir failed: %w", err)
		}
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("open destination failed: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}
	return nil
}

func Cwd() (string, error) { return os.Getwd() }

func HomeDir() (string, error) { return os.UserHomeDir() }

func TempDir() string { return os.TempDir() }
