package plugins

import "cmds/sugar"

func LaunchApp(cmd string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func LaunchChrome() {
	LaunchApp("chrome --proxy-server=socks5://127.0.0.1:7891")
}

func LaunchEdge() {
	LaunchApp("edge --proxy-server=socks5://127.0.0.1:7891")
}
