package plugins

import (
	"fmt"
	"time"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

func GetHostName() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd := "hostname"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := stdout
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func GetIP() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd := "hostname -i"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := stdout
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func GetCurrentDatetime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := time.Now().Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func GetCurrentUnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	content := fmt.Sprintf("%d", time.Now().Unix())
	sugar.Notify(fmt.Sprintf("get success: %s", content))
	changed := clipboard.Write(clipboard.FmtText, []byte(content))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}
