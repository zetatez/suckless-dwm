package svc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"assistant/internal/bootstrap/psl"
)

type fileLocationPattern struct {
	re       *regexp.Regexp
	fileIdx  int
	lineIdx  int
	colIdx   int
	trimTail func(string) string
}

var fileLocationPatterns = []fileLocationPattern{
	{
		regexp.MustCompile(`(?m)(~[A-Za-z0-9_./\-.~]+)`),
		1,
		0,
		0,
		func(s string) string { return s },
	}, // ~/path
	{
		regexp.MustCompile(`(?m)(/[^:\s]+):(\d+)(?::(\d+))?`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // /path:line:col
	{
		regexp.MustCompile(`(?m)\s+at\s+(/[^:\s]+):(\d+):(\d+)`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // at /path:line:col
	{
		regexp.MustCompile(`(?m)File\s+"([^"]+)",\s+line\s+(\d+)`),
		1,
		2,
		0,
		func(s string) string { return s },
	}, // File "path", line N
	{
		regexp.MustCompile(`(?m)-->\s+(/[^:\s]+):(\d+):(\d+)`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // --> /path:line:col
	{
		regexp.MustCompile(`(?m)\((/[^:()]+):(\d+):(\d+)\)`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // (path:line:col)
	{
		regexp.MustCompile(`(?m)(~[^:\s]+):(\d+)(?::(\d+))?`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // ~/path:line:col
	{
		regexp.MustCompile(`(?m)([A-Za-z0-9_./\-~]+\.[A-Za-z0-9]+):(\d+)(?::(\d+))?`),
		1,
		2,
		3,
		func(s string) string { return strings.TrimRight(strings.TrimRight(s, ")"), ":") },
	}, // file:line:col
}

func (s *Service) extractFileLocation(text string) (file string, line int, col int, ok bool) {
	for _, p := range fileLocationPatterns {
		m := p.re.FindStringSubmatch(text)
		if len(m) == 0 {
			continue
		}

		candidate := p.trimTail(strings.TrimSpace(m[p.fileIdx]))

		if v, _ := strconv.Atoi(m[p.lineIdx]); v > 0 {
			line = v
		}
		if p.colIdx > 0 {
			if v, _ := strconv.Atoi(m[p.colIdx]); v >= 0 {
				col = v
			}
		}

		if strings.HasPrefix(candidate, "~") {
			candidate = filepath.Join(os.Getenv("HOME"), candidate[1:])
		}

		if filepath.IsAbs(candidate) {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, line, col, true
			}
			continue
		}

		if abs, err := filepath.Abs(candidate); err == nil {
			candidate = abs
		} else if out, _, err := runScript("bash", fmt.Sprintf("fd -H -1 -f '%s'", candidate)); err == nil && out != "" {
			candidate = strings.TrimSpace(out)
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate, line, col, true
		}
	}
	return "", 0, 0, false
}

func (s *Service) extractURL(text string) (string, bool) {
	m := regexp.MustCompile(`https?://[^\s)]+`).FindString(strings.TrimSpace(text))
	if m != "" {
		return m, true
	}
	return "", false
}

func (s *Service) isAbsFileAndExist(text string) bool {
	if strings.HasPrefix(text, "~") {
		text = filepath.Join(os.Getenv("HOME"), text[1:])
	}
	if !filepath.IsAbs(text) {
		return false
	}
	_, err := os.Stat(text)
	return err == nil
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

	if url, ok := s.extractURL(text); ok {
		err = s.OpenURL("chrome", url)
		return fmt.Sprintf("opened URL: %s", url), err
	}

	term := psl.GetConfig().Svc.DefaultTerminal

	if s.isAbsFileAndExist(text) {
		_, _, err = runScript("bash", fmt.Sprintf("%s -e nvim %s", term, text))
		return fmt.Sprintf("opened %s", text), err
	}

	if file, line, col, ok := s.extractFileLocation(text); ok {
		var vimcmd string
		switch {
		case line > 0 && col > 0:
			vimcmd = fmt.Sprintf("call cursor(%d,%d)", line, col)
			_, _, err = runScript("bash", fmt.Sprintf("%s -e nvim +'%s' %s", term, vimcmd, file))
		case line > 0:
			_, _, err = runScript("bash", fmt.Sprintf("%s -e nvim +%d %s", term, line, file))
		default:
			_, _, err = runScript("bash", fmt.Sprintf("%s -e nvim %s", term, file))
		}
		msg := fmt.Sprintf("opened %s", file)
		return msg, err
	}

	err = s.SearchWeb(text)
	return fmt.Sprintf("searched: %s", text), err
}
