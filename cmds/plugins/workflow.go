package plugins

import (
	"fmt"
	"strings"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

func HandleCopied() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	content := strings.TrimSpace(string(text))
	switch {
	case sugar.Exists(content) && sugar.IsFile(content):
		sugar.Lazy("open", content)
		return
	case sugar.IsUrl(content):
		OpenWeb(content)()
		return
	default:
		SearchFromWeb(content)
	}
}

func UmountXYZ() {
	cmd := "echo '/x\n/y\n/z'|dmenu -p 'umount'"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	choice := stdout
	cmd = fmt.Sprintf("sudo umount %s", choice)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("umount success")
}

func WifiConnect() {
	cmd := "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to wifi'", stdout)
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	essid := strings.TrimSpace(stdout)
	if essid == "" {
		return
	}
	cmd = "dmenu < /dev/null -p 'password'"
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	password := strings.TrimSpace(stdout)
	cmd = fmt.Sprintf("nmcli device wifi connect %s password %s", essid, password)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("wifi connect success")
}
