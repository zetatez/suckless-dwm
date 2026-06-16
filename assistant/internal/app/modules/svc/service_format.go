package svc

import (
	"encoding/json"
	"fmt"
	"go/format"
	"strings"

	"gopkg.in/yaml.v3"
)

func (s *Service) Format(language string) (string, error) {
	content, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	if content == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	var result string

	switch language {
	case "json":
		var doc any
		if err := json.Unmarshal([]byte(content), &doc); err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
		formatted, e := json.MarshalIndent(doc, "", "  ")
		if e != nil {
			return "", fmt.Errorf("format JSON failed: %w", e)
		}
		result = string(formatted)
	case "yml", "yaml":
		var doc any
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return "", fmt.Errorf("invalid YAML: %w", err)
		}
		formatted, e := yaml.Marshal(&doc)
		if e != nil {
			return "", fmt.Errorf("format YAML failed: %w", e)
		}
		result = string(formatted)
	case "sql":
		stdout, stderr, e := runScript("python", fmt.Sprintf(
			`import sqlparse; print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='lower'))`, content))
		if e != nil {
			return "", fmt.Errorf("format SQL failed: %s", stderr)
		}
		result = strings.TrimSpace(stdout)
	case "go":
		formatted, e := format.Source([]byte(content))
		if e != nil {
			return "", fmt.Errorf("format Go failed: %w", e)
		}
		result = string(formatted)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	if _, err := s.pushClipboard(result, fmt.Sprintf("format %s success", language)); err != nil {
		return result, err
	}
	return result, nil
}
