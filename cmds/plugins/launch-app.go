package plugins

import "cmds/sugar"

func LaunchApp(cmd string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(cmd)
	}
}
