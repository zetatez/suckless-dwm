package plugins

import (
	"fmt"

	"cmds/sugar"
)

func SearchFromWeb(content string) {
	sugar.NewExecService().RunScriptShell(
		fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 https://cn.bing.com/search?q='%s'", content),
	)
}
