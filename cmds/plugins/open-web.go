package plugins

import (
	"fmt"

	"cmds/sugar"
)

func OpenWeb(url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"chrome --proxy-server=%s %s",
				ProxyServer,
				url,
			),
		)
	}
}
