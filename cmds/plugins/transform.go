package plugins

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

func TransformDatetime2UnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	t, err := time.Parse(time.DateTime, strings.TrimSpace(string(text)))
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := fmt.Sprintf("%d", t.Unix())
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func TransformUnixSec2DateTime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	unix, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		sugar.Notify(err)
		return
	}
	datetime := time.Unix(unix, 0).Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", datetime))
	changed := clipboard.Write(clipboard.FmtText, []byte(datetime))
	<-changed
	sugar.Notify("previous clipboard expired")
}
