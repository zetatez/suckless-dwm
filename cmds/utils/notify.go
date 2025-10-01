package utils

import (
	"fmt"
)

func Notify(msg ...any) {
	message := fmt.Sprint(msg...)
	CmdMap := map[string][]string{
		"linux":  {"bash", fmt.Sprintf("notify-send '%s'", message)},
		"darwin": {"bash", fmt.Sprintf(`osascript -e 'display notification "%s" with title "%s"'`, message, "msg")},
	}
	if cmd, ok := CmdMap[GetOSType()]; ok {
		RunScript(cmd[0], cmd[1])
		return
	}
	fmt.Println("Unsupported OS, exiting...")
	return
}
