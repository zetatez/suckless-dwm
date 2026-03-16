package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func appleScriptQuote(s string) string {
	// AppleScript string literal uses double quotes.
	esc := strings.ReplaceAll(s, "\\", "\\\\")
	esc = strings.ReplaceAll(esc, "\"", "\\\"")
	return "\"" + esc + "\""
}

func Notify(msg ...any) {
	message := fmt.Sprint(msg...)
	osType := GetOSType()

	switch osType {
	case "linux":
		_ = exec.Command("notify-send", message).Run()
	case "darwin":
		cmd := fmt.Sprintf(
			"display notification %s with title %s",
			appleScriptQuote(message),
			appleScriptQuote("msg"),
		)
		_ = exec.Command("osascript", "-e", cmd).Run()
	default:
		fmt.Println(message)
	}
}
