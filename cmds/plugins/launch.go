package plugins

import (
	"fmt"

	"cmds/utils"
)

func LaunchApp(cmd string) func() {
	return func() {
		utils.RunScript("bash", cmd)
	}
}

func LaunchChrome() {
	LaunchApp(
		fmt.Sprintf("chrome --proxy-server=%s --new-window", ProxyServer),
	)()
}

func LaunchQutebrowser() {
	LaunchApp(
		fmt.Sprintf("qutebrowser --set content.proxy '%s'", ProxyServer),
	)()
}

func LaunchFileManager() {
	LaunchApp("thunar ~")()
}

func LaunchInkscape() {
	LaunchApp("inkscape")()
}

func LaunchKrita() {
	LaunchApp("krita")()
}

func LaunchObsidian() {
	LaunchApp("obsidian")()
}

func LaunchSublime() {
	LaunchApp("subl")()
}

func LaunchXournal() {
	LaunchApp("xournalpp")()
}
