package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type GrepResult struct {
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Content   string `json:"content"`
	MatchInfo string `json:"match_info"`
}

type Grep struct {
	pattern        string
	paths          []string
	invertMatch    bool
	context        int
	recursive      bool
	filePattern    *regexp.Regexp
	filePatternErr error
	countOnly      bool
}

func NewGrep(pattern string, paths ...string) *Grep {
	return &Grep{pattern: pattern, paths: paths, recursive: true}
}

func (g *Grep) WithInvertMatch() *Grep {
	g.invertMatch = true
	return g
}

func (g *Grep) WithContext(n int) *Grep {
	g.context = n
	return g
}

func (g *Grep) WithRecursive(b bool) *Grep {
	g.recursive = b
	return g
}

func (g *Grep) WithFilePattern(pattern string) *Grep {
	re, err := regexp.Compile(pattern)
	g.filePattern = re
	g.filePatternErr = err
	return g
}

func (g *Grep) WithCountOnly() *Grep {
	g.countOnly = true
	return g
}

func (g *Grep) Exec() ([]GrepResult, error) {
	if g.filePatternErr != nil {
		return nil, fmt.Errorf("invalid file pattern: %w", g.filePatternErr)
	}
	if _, err := regexp.Compile(g.pattern); err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	// Fast path: use ripgrep if available.
	if _, ok := lookPath("rg"); ok {
		if res, ok, err := g.execRG(); ok {
			return res, err
		}
	}

	// Fallback: pure Go walk+scan.
	re, _ := regexp.Compile(g.pattern)

	var results []GrepResult
	for _, root := range g.paths {
		err := filepath.Walk(root, func(filePath string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if info.IsDir() {
				if !g.recursive && filePath != root {
					return filepath.SkipDir
				}
				return nil
			}
			if g.filePattern != nil && !g.filePattern.MatchString(info.Name()) {
				return nil
			}

			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("open file failed (%s): %w", filePath, err)
			}
			defer f.Close()

			scanner := NewLineScanner(f)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				if g.invertMatch {
					if !re.MatchString(line) {
						results = append(results, GrepResult{Path: filePath, Line: lineNum, Content: line, MatchInfo: re.String()})
					}
					continue
				}

				loc := re.FindStringIndex(line)
				if loc == nil {
					continue
				}

				match := line[loc[0]:loc[1]]
				ctxLine := line
				if g.context > 0 {
					start := max(0, loc[0]-g.context)
					end := min(len(line), loc[1]+g.context)
					ctxLine = line[start:end]
				}
				results = append(results, GrepResult{
					Path:      filePath,
					Line:      lineNum,
					Content:   ctxLine,
					MatchInfo: fmt.Sprintf("'%s' at pos %d", match, loc[0]),
				})
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan file failed (%s): %w", filePath, err)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk path failed (%s): %w", root, err)
		}
	}

	if g.countOnly {
		return []GrepResult{{Path: "(count)", Line: len(results)}}, nil
	}
	return results, nil
}

func (g *Grep) execRG() (results []GrepResult, ok bool, err error) {
	// Keep behavior consistent with the pure-Go fallback.
	if g.invertMatch {
		return nil, false, nil
	}

	args := []string{"--json", "--no-messages", "--hidden", "--no-ignore", "--no-ignore-vcs"}
	if !g.recursive {
		args = append(args, "--max-depth", "1")
	}

	if g.filePattern != nil {
		if glob, convOK := regexNameToRGGlob(g.filePattern.String()); convOK {
			args = append(args, "--glob", glob)
		} else {
			return nil, false, nil
		}
	}

	args = append(args, g.pattern)
	if len(g.paths) == 0 {
		args = append(args, ".")
	} else {
		args = append(args, g.paths...)
	}

	count := 0
	stderr, runErr := runRipgrepJSON(args, func(md rgMatchData) error {
		count++
		if g.countOnly {
			return nil
		}

		line := strings.TrimRight(md.Lines.Text, "\r\n")
		start, end := 0, 0
		matchText := ""
		if len(md.Submatches) > 0 {
			start = md.Submatches[0].Start
			end = md.Submatches[0].End
			if start >= 0 && end >= start && end <= len(line) {
				matchText = line[start:end]
			}
		}

		ctxLine := line
		if g.context > 0 && len(md.Submatches) > 0 {
			cs := start - g.context
			ce := end + g.context
			if cs < 0 {
				cs = 0
			}
			if ce > len(line) {
				ce = len(line)
			}
			if cs <= ce {
				ctxLine = line[cs:ce]
			}
		}

		mi := fmt.Sprintf("pos %d", start)
		if matchText != "" {
			mi = fmt.Sprintf("'%s' at pos %d", matchText, start)
		}

		results = append(results, GrepResult{
			Path:      md.Path.Text,
			Line:      md.LineNumber,
			Content:   ctxLine,
			MatchInfo: mi,
		})
		return nil
	})
	if runErr != nil {
		// If rg itself is available but fails for some reason, fall back to pure-Go.
		if strings.TrimSpace(stderr) != "" {
			return nil, false, nil
		}
		return nil, false, nil
	}

	if g.countOnly {
		return []GrepResult{{Path: "(count)", Line: count}}, true, nil
	}
	return results, true, nil
}

func lookPath(cmd string) (string, bool) {
	p, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}
	return p, true
}

type rgEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type rgMatchData struct {
	Path struct {
		Text string `json:"text"`
	} `json:"path"`
	Lines struct {
		Text string `json:"text"`
	} `json:"lines"`
	LineNumber int `json:"line_number"`
	Submatches []struct {
		Match struct {
			Text string `json:"text"`
		} `json:"match"`
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"submatches"`
}

func runRipgrepJSON(args []string, onMatch func(md rgMatchData) error) (string, error) {
	rgPath, ok := lookPath("rg")
	if !ok {
		return "", fmt.Errorf("rg not found")
	}

	cmd := exec.Command(rgPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("rg stdout pipe: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("rg start: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 128*1024), 8*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var ev rgEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			return stderr.String(), fmt.Errorf("rg json decode: %w", err)
		}
		if ev.Type != "match" {
			continue
		}
		var md rgMatchData
		if err := json.Unmarshal(ev.Data, &md); err != nil {
			return stderr.String(), fmt.Errorf("rg match decode: %w", err)
		}
		if err := onMatch(md); err != nil {
			return stderr.String(), err
		}
	}
	if err := scanner.Err(); err != nil {
		return stderr.String(), fmt.Errorf("rg scan: %w", err)
	}
	if err := stdout.Close(); err != nil {
		return stderr.String(), fmt.Errorf("rg stdout close: %w", err)
	}

	err = cmd.Wait()
	if err == nil {
		return stderr.String(), nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() == 1 {
			return stderr.String(), nil
		}
		return stderr.String(), fmt.Errorf("rg failed (exit %d): %s", exitErr.ExitCode(), strings.TrimSpace(stderr.String()))
	}
	return stderr.String(), fmt.Errorf("rg wait: %w", err)
}

func regexNameToRGGlob(pattern string) (string, bool) {
	p := strings.TrimSpace(pattern)
	p = strings.TrimPrefix(p, "^")
	p = strings.TrimSuffix(p, "$")

	if strings.Contains(p, "/") || strings.Contains(p, `\\`) {
		return "", false
	}

	for _, pref := range []string{".*", ".+", "(.*)", "(.+)"} {
		if strings.HasPrefix(p, pref) {
			p = strings.TrimPrefix(p, pref)
			break
		}
	}
	if !strings.HasPrefix(p, `\\.`) {
		return "", false
	}
	p = strings.TrimPrefix(p, `\\.`)
	if p == "" {
		return "", false
	}

	if strings.HasPrefix(p, "(") && strings.HasSuffix(p, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(p, "("), ")")
		parts := strings.Split(inner, "|")
		for _, part := range parts {
			if part == "" || strings.ContainsAny(part, `.*+?[]{}()\\`) {
				return "", false
			}
		}
		return "*.{" + strings.Join(parts, ",") + "}", true
	}

	if strings.ContainsAny(p, `.*+?[]{}()|\\`) {
		return "", false
	}
	return "*." + p, true
}
