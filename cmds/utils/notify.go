package utils

import "fmt"

func Notify(msg ...any) {
	message := fmt.Sprint(msg...)
	osType := GetOSType()

	switch osType {
	case "linux":
		RunScript("bash", fmt.Sprintf("notify-send '%s'", message))
	case "darwin":
		RunScript("bash", fmt.Sprintf(`osascript -e 'display notification "%s" with title "%s"'`, message, "msg"))
	default:
		fmt.Println("Unsupported OS, exiting...")
	}
}
