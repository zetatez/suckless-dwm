package plugins

import (
	"fmt"

	"cmds/sugar"
)

func LaunchApp(cmd string) func() {
	return func() {
		sugar.NewExecService().RunScript("bash", cmd)
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

func LaunchQuteBrowser() {
	LaunchApp(
		fmt.Sprintf(
			"qutebrowser --set content.proxy %s --target window",
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
