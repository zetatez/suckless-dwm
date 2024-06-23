package plugins

import (
	"fmt"

	"cmds/sugar"
)

func SearchFromWeb(content string) {
	// sugar.NewExecService().RunScriptShell(
	// 	fmt.Sprintf(
	// 		"chrome --proxy-server=%s https://cn.bing.com/search?q='%s'",
	// 		ProxyServer,
	// 		content,
	// 	),
	// )
	sugar.NewExecService().RunScriptShell(
		fmt.Sprintf(
			"chrome --proxy-server=%s https://www.google.com/search?q='%s'",
			ProxyServer,
			content,
		),
	)
}
