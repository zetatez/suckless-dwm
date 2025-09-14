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
		fmt.Sprintf(
			"chrome --proxy-server=%s --new-window",
			ProxyServer,
		),
	)()
}

func LaunchEdge() {
	LaunchApp(
		fmt.Sprintf(
			"edge --proxy-server=%s --new-window",
			ProxyServer,
		),
	)()
}

func LaunchThunar() {
	LaunchApp("thunar ~")()
}
