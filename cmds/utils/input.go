package utils

import (
	"fmt"
	"strings"
)

func GetInput(prompt string) (input string, err error) {
	script := fmt.Sprintf(
		"rofi -dmenu < /dev/null -p '%s'",
		prompt,
	)
	stdout, _, err := RunScript("bash", script)
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(stdout)
	return input, nil
}
