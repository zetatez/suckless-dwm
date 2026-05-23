package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

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

// runRipgrepLines runs rg and calls onLine for each stdout line.
// It treats exit code 1 (no matches) as success.
func runRipgrepLines(args []string, onLine func(line string) error) (string, error) {
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
		if err := onLine(scanner.Text()); err != nil {
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

// runRipgrepJSON runs `rg --json` and calls onMatch for each match event.
// It treats exit code 1 (no matches) as success.
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
	// JSON lines can be long (very long paths / long matched lines).
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
		// Exit code 1 means no matches.
		if exitErr.ExitCode() == 1 {
			return stderr.String(), nil
		}
		return stderr.String(), fmt.Errorf("rg failed (exit %d): %s", exitErr.ExitCode(), strings.TrimSpace(stderr.String()))
	}
	return stderr.String(), fmt.Errorf("rg wait: %w", err)
}

func regexNameToRGGlob(pattern string) (string, bool) {
	// Best-effort conversion for common "extension" regexes:
	//   \\.go$            -> *.go
	//   .*\\.(go|ts)$      -> *.{go,ts}
	//   ^.*\\.(go|ts)$     -> *.{go,ts}
	//   ^.+\\.go$          -> *.go
	// Anything else: return false.

	p := strings.TrimSpace(pattern)
	p = strings.TrimPrefix(p, "^")
	p = strings.TrimSuffix(p, "$")

	// Strip leading "(.*/)?"-style patterns; Find/Grep filePattern is applied to base name.
	if strings.Contains(p, "/") || strings.Contains(p, `\\`) {
		return "", false
	}

	// Accept a few common prefixes.
	for _, pref := range []string{".*", ".+", "(.*)", "(.+)"} {
		if strings.HasPrefix(p, pref) {
			p = strings.TrimPrefix(p, pref)
			break
		}
	}
	// Require literal dot extension.
	if !strings.HasPrefix(p, `\\.`) {
		return "", false
	}
	p = strings.TrimPrefix(p, `\\.`)
	if p == "" {
		return "", false
	}

	// (a|b|c)
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

	// Single extension.
	if strings.ContainsAny(p, `.*+?[]{}()|\\`) {
		return "", false
	}
	return "*." + p, true
}

func baseName(path string) string {
	// filepath.Base treats trailing slashes specially; rg --json path.text should be file path.
	return filepath.Base(path)
}
