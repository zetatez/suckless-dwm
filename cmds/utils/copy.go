package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyDir(src string, dst string) (err error) {
	if strings.TrimSpace(src) == strings.TrimSpace(dst) {
		return fmt.Errorf("src path %s is equal to dst path %s", src, dst)
	}

	if !IsDirExists(src) {
		return fmt.Errorf("src path %s is not exist", src)
	}

	if !IsDirExists(dst) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
	}

	absSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	err = filepath.Walk(
		absSrc,
		func(s string, info os.FileInfo, err error) error {
			if s == absSrc {
				return nil
			}
			if info == nil {
				return err
			}
			d := strings.ReplaceAll(s, absSrc, absDst)
			if info.IsDir() {
				if !IsDirExists(d) {
					if err = os.MkdirAll(d, os.ModePerm); err != nil {
						return err
					}
				}
			} else {
				err = CopyFile(s, d)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	return err
}

func CopyFile(src string, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	sinfo, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, sinfo.Mode())
	if err != nil {
		return err
	}
	return nil
}
