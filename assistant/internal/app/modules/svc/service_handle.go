package svc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var fileLocationPatterns = []struct {
	re   *regexp.Regexp
	file int
	line int
	col  int
}{
	{regexp.MustCompile(`(?m)(/[^:\s]+):(\d+)(?::(\d+))?`), 1, 2, 3},
	{regexp.MustCompile(`(?m)([A-Za-z0-9_./\-~]+\.[A-Za-z0-9]+):(\d+)(?::(\d+))?`), 1, 2, 3},
	{regexp.MustCompile(`(?m)File\s+"([^"]+)",\s+line\s+(\d+)`), 1, 2, 0},
	{regexp.MustCompile(`(?m)\((/[^:()]+):(\d+):(\d+)\)`), 1, 2, 3},
	{regexp.MustCompile(`(?m)\s+at\s+(/[^:\s]+):(\d+):(\d+)`), 1, 2, 3},
	{regexp.MustCompile(`(?m)-->\s+(/[^:\s]+):(\d+):(\d+)`), 1, 2, 3},
}

func (s *Service) extractFileLocation(text string) (file string, line, col int, ok bool) {
	for _, p := range fileLocationPatterns {
		m := p.re.FindStringSubmatch(text)
		if len(m) == 0 {
			continue
		}
		candidate := strings.TrimSpace(m[p.file])
		candidate = strings.TrimSuffix(candidate, ")")
		candidate = strings.TrimSuffix(candidate, ":")
		l, err := strconv.Atoi(m[p.line])
		if err != nil || l <= 0 {
			continue
		}
		c := 0
		if p.col > 0 && p.col < len(m) && m[p.col] != "" {
			if x, e := strconv.Atoi(m[p.col]); e == nil {
				c = x
			}
		}
		if !filepath.IsAbs(candidate) {
			if abs, e := filepath.Abs(candidate); e == nil {
				candidate = abs
			}
		}
		if _, e := os.Stat(candidate); e == nil {
			return candidate, l, c, true
		}
	}
	return "", 0, 0, false
}

func (s *Service) extractMarkdownURL(text string) (string, bool) {
	m := regexp.MustCompile(`\[[^\]]*\]\((https?://[^\s)]+)\)`).FindStringSubmatch(strings.TrimSpace(text))
	if len(m) == 2 {
		return m[1], true
	}
	return "", false
}

func isURL(text string) bool {
	return strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://")
}

func (s *Service) existsAndIsFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func (s *Service) HandleClipboard() (string, error) {
	text, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	if file, line, col, ok := s.extractFileLocation(text); ok {
		term := "st"
		cmd := fmt.Sprintf("%s -e nvim +%d %s", term, line, file)
		if col > 0 {
			cmd = fmt.Sprintf("%s -e nvim +'%s' %s", term, fmt.Sprintf("call cursor(%d,%d)", line, col), file)
		}
		_, _, err := runScript("bash", cmd)
		return fmt.Sprintf("opened %s:%d:%d", file, line, col), err
	}

	if s.existsAndIsFile(text) {
		err := startScript("bash", fmt.Sprintf("st -e yazi '%s'", text))
		return fmt.Sprintf("opened file: %s", text), err
	}

	if url, ok := s.extractMarkdownURL(text); ok {
		err := s.OpenURL("chrome", url)
		return fmt.Sprintf("opened URL: %s", url), err
	}

	if isURL(text) {
		err := s.OpenURL("chrome", text)
		return fmt.Sprintf("opened URL: %s", text), err
	}

	err = s.SearchWeb(text)
	return fmt.Sprintf("searched: %s", text), err
}
