package plugins

import (
	"fmt"
	"regexp"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

func JumpToCodeFromLog() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	textbyte := clipboard.Read(clipboard.FmtText)
	text := string(textbyte)
	regex := `(?P<filepath>/[^\:]+):(?P<row>\d+)\s+`
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(text)
	if len(match) < 3 {
		sugar.Notify("not match")
		return
	}
	filepath := match[1]
	row := match[2]
	cmd := fmt.Sprintf(
		"st -e nvim +%s %s",
		row,
		filepath,
	)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}
