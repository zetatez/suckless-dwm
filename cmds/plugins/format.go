package plugins

import (
	"encoding/json"
	"fmt"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

func ReturnFormatJson() func() {
	return FormatJson
}

func FormatJson() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	doc := map[string]interface{}{}
	err = json.Unmarshal(text, &doc)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify(fmt.Sprintf("format json success: \n%s", string(formatedText)))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	select {
	case <-changed:
		sugar.Notify("previous formated json expired")
	}
}

func ReturnFormatSql() func() {
	return FormatSql
}

func FormatSql() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	cmd := `
import sqlparse
print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='upper'))
`
	cmd = fmt.Sprintf(cmd, string(text))
	stdout, _, err := sugar.NewExecService().RunScriptPython(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := stdout
	sugar.Notify(fmt.Sprintf("format sql success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous formated json expired")
	}
}
