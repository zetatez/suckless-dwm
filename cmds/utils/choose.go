package utils

import (
	"fmt"
	"strings"
)

func Choose(prompt string, list []string) (item string, err error) {
	script := fmt.Sprintf(
		"echo '%s'|rofi -dmenu -p '%s'",
		strings.Join(list, "\n"),
		prompt,
	)
	stdout, _, err := RunScript("bash", script)
	if err != nil {
		return "", err
	}
	item = strings.TrimSpace(stdout)
	return item, nil
}
