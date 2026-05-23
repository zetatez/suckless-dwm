package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type GlobResult struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

type Glob struct {
	baseDir     string
	pattern     string
	includeFile bool
	includeDir  bool
	maxResults  int
}

func NewGlob(pattern string) *Glob {
	return &Glob{baseDir: ".", pattern: pattern, includeFile: true, includeDir: true}
}

func (g *Glob) InDir(baseDir string) *Glob {
	if strings.TrimSpace(baseDir) != "" {
		g.baseDir = baseDir
	}
	return g
}

func (g *Glob) WithFilesOnly() *Glob {
	g.includeFile = true
	g.includeDir = false
	return g
}

func (g *Glob) WithDirsOnly() *Glob {
	g.includeFile = false
	g.includeDir = true
	return g
}

func (g *Glob) WithMaxResults(n int) *Glob {
	g.maxResults = n
	return g
}

func (g *Glob) Exec() ([]GlobResult, error) {
	pat := strings.TrimSpace(g.pattern)
	if pat == "" {
		return nil, fmt.Errorf("%w: empty pattern", ErrInvalidPath)
	}
	if filepath.IsAbs(pat) {
		return nil, fmt.Errorf("%w: absolute pattern not allowed: %s", ErrInvalidPath, pat)
	}
	// Reject explicit '..' traversal in patterns.
	if c := filepath.Clean(pat); strings.HasPrefix(c, ".."+string(filepath.Separator)) || c == ".." {
		return nil, fmt.Errorf("%w: pattern escapes base dir: %s", ErrInvalidPath, pat)
	}

	base := g.baseDir
	if strings.TrimSpace(base) == "" {
		base = "."
	}

	patternSlash := filepath.ToSlash(strings.TrimPrefix(pat, "./"))

	var out []GlobResult
	stopErr := fmt.Errorf("stop glob walk")
	err := filepath.Walk(base, func(full string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(base, full)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(strings.TrimPrefix(rel, "./"))
		if relSlash == "." {
			return nil
		}

		matched, err := matchGlob(patternSlash, relSlash)
		if err != nil {
			return err
		}
		if !matched {
			return nil
		}

		if info.IsDir() {
			if !g.includeDir {
				return nil
			}
		} else {
			if !g.includeFile {
				return nil
			}
		}

		out = append(out, GlobResult{Path: full, IsDir: info.IsDir()})
		if g.maxResults > 0 && len(out) >= g.maxResults {
			return stopErr
		}
		return nil
	})
	if err != nil {
		if err.Error() != stopErr.Error() {
			return nil, fmt.Errorf("glob failed: %w", err)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

// matchGlob matches a slash-separated path against a glob pattern.
// It supports ** to match zero or more path segments.
func matchGlob(patternSlash, relSlash string) (bool, error) {
	pp := strings.Split(patternSlash, "/")
	rp := strings.Split(relSlash, "/")

	var rec func(pi, ri int) (bool, error)
	rec = func(pi, ri int) (bool, error) {
		if pi == len(pp) {
			return ri == len(rp), nil
		}
		seg := pp[pi]
		if seg == "**" {
			for k := ri; k <= len(rp); k++ {
				ok, err := rec(pi+1, k)
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
			}
			return false, nil
		}
		if ri >= len(rp) {
			return false, nil
		}
		ok, err := path.Match(seg, rp[ri])
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		return rec(pi+1, ri+1)
	}

	return rec(0, 0)
}
