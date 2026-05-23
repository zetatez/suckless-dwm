package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FindResult struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
	Content string `json:"content,omitempty"`
}

type Find struct {
	root           string
	namePattern    *regexp.Regexp
	namePatternErr error
	typeFilter     string
	minSize        int64
	maxSize        int64
	searchHidden   bool
	withContent    bool
}

func FindNew(root string) *Find {
	return &Find{root: root, searchHidden: true}
}

func (f *Find) WithNamePattern(pattern string) *Find {
	re, err := regexp.Compile(pattern)
	f.namePattern = re
	f.namePatternErr = err
	return f
}

func (f *Find) WithTypeFilter(t string) *Find {
	f.typeFilter = t
	return f
}

func (f *Find) WithMinSize(bytes int64) *Find {
	f.minSize = bytes
	return f
}

func (f *Find) WithMaxSize(bytes int64) *Find {
	f.maxSize = bytes
	return f
}

func (f *Find) WithSearchHidden() *Find {
	f.searchHidden = true
	return f
}

func (f *Find) WithContent() *Find {
	f.withContent = true
	return f
}

func (f *Find) Exec() ([]FindResult, error) {
	if f.namePatternErr != nil {
		return nil, fmt.Errorf("invalid name pattern: %w", f.namePatternErr)
	}

	var results []FindResult
	err := filepath.Walk(f.root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() && !f.searchHidden && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if f.namePattern != nil && !f.namePattern.MatchString(info.Name()) {
			return nil
		}

		if f.typeFilter != "" {
			switch f.typeFilter {
			case "f":
				if info.IsDir() {
					return nil
				}
			case "d":
				if !info.IsDir() {
					return nil
				}
			default:
				return fmt.Errorf("invalid type filter: %q", f.typeFilter)
			}
		}

		if f.minSize > 0 && info.Size() < f.minSize {
			return nil
		}
		if f.maxSize > 0 && info.Size() > f.maxSize {
			return nil
		}

		fi := FindResult{
			Path:    path,
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}

		if f.withContent && !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file failed (%s): %w", path, err)
			}
			fi.Content = string(content)
		}

		results = append(results, fi)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	return results, nil
}
