package plugins

import (
	"encoding/json"
	"fmt"

	"cmds/sugar"

	"golang.design/x/clipboard"
	"gopkg.in/yaml.v3"
)

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
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	<-changed
	sugar.Notify("previous clipboard expired")
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
print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='lower'))
`
	cmd = fmt.Sprintf(cmd, string(text))
	stdout, _, err := sugar.NewExecService().RunScript("python", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := stdout
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	<-changed
	sugar.Notify("previous clipboard expired")
}

func FormatYaml() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	doc := map[interface{}]interface{}{}
	err = yaml.Unmarshal(text, &doc)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText, err := yaml.Marshal(&doc)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	<-changed
	sugar.Notify("previous clipboard expired")
}
