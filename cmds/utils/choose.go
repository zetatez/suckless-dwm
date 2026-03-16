package utils

import (
	"fmt"
	"strings"
)

func Choose(prompt string, list []string) (item string, err error) {
	content := strings.Join(list, "\n")
	script := fmt.Sprintf("printf %s | rofi -dmenu -p %s", ShellSingleQuote(content), ShellSingleQuote(prompt))
	stdout, _, err := RunScript("bash", script)
	if err != nil {
		return "", err
	}
	item = strings.TrimSpace(stdout)
	return item, nil
}
