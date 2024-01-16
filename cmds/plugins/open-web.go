package plugins

import (
	"fmt"

	"cmds/sugar"
)

func OpenWeb(url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 %s", url),
		)
	}
}
