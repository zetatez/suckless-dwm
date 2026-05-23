package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ToJSON(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}
	return string(data), nil
}

func FromJSON(jsonStr string, v any) error {
	if err := json.Unmarshal([]byte(jsonStr), v); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}
	return nil
}

func ToJSONCompact(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}
	return string(data), nil
}

func ValidateJSON(jsonStr string) error {
	var tmp any
	if err := json.Unmarshal([]byte(jsonStr), &tmp); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return nil
}

func PrettyJSON(jsonStr string, indent string) (string, error) {
	if indent == "" {
		indent = "  "
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(jsonStr), "", indent); err != nil {
		return "", fmt.Errorf("indent failed: %w", err)
	}
	return buf.String(), nil
}

// ExtractJSONObject returns the first valid JSON object substring found in input.
// It is string-aware (ignores braces inside JSON strings).
func ExtractJSONObject(input string) (string, error) {
	start := -1
	depth := 0
	inString := false
	escapeNext := false

	for i, r := range input {
		if escapeNext {
			escapeNext = false
			continue
		}

		if r == '\\' && inString {
			escapeNext = true
			continue
		}

		if r == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if r == '{' {
			if depth == 0 {
				start = i
			}
			depth++
			continue
		}
		if r == '}' {
			depth--
			if depth == 0 && start != -1 {
				return input[start : i+1], nil
			}
		}
	}

	return "{}", errors.New("no valid json object found")
}

// ExtractJSONArray returns the first valid JSON array substring found in input.
// It is string-aware (ignores brackets inside JSON strings).
func ExtractJSONArray(input string) (string, error) {
	start := -1
	depth := 0
	inString := false
	escapeNext := false

	for i, r := range input {
		if escapeNext {
			escapeNext = false
			continue
		}

		if r == '\\' && inString {
			escapeNext = true
			continue
		}

		if r == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if r == '[' {
			if depth == 0 {
				start = i
			}
			depth++
			continue
		}
		if r == ']' {
			depth--
			if depth == 0 && start != -1 {
				return input[start : i+1], nil
			}
		}
	}

	return "[]", errors.New("no valid json array found")
}

// ExtractAndValidateJSONObject extracts the first JSON object from input and
// validates it parses as JSON.
func ExtractAndValidateJSONObject(input string) (string, error) {
	obj, err := ExtractJSONObject(input)
	if err != nil {
		return "{}", err
	}
	var tmp any
	if err := json.Unmarshal([]byte(obj), &tmp); err != nil {
		return "{}", err
	}
	return obj, nil
}

// writeFileAtomic writes data to a temporary file and atomically renames it
// to the target path, ensuring the file is never in a partial state.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	if err := os.Chmod(path, perm); err != nil {
		return fmt.Errorf("chmod file: %w", err)
	}
	return nil
}

// CleanJSONResponse removes surrounding markdown code fences (``` or ```json)
// from a model response.
func CleanJSONResponse(content string) string {
	content = strings.TrimSpace(content)

	markdownStart := strings.Index(content, "```json")
	codeStart := strings.Index(content, "```")

	if markdownStart != -1 && (codeStart == -1 || markdownStart < codeStart) {
		end := strings.LastIndex(content, "```")
		if end > markdownStart {
			content = strings.TrimSpace(content[markdownStart+7 : end])
		}
	} else if codeStart != -1 {
		end := strings.LastIndex(content, "```")
		if end > codeStart {
			lines := strings.Split(content[codeStart+3:end], "\n")
			if len(lines) > 1 {
				content = strings.TrimSpace(strings.Join(lines[1:], "\n"))
			}
		}
	}

	return content
}
