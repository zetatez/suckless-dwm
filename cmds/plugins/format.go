package plugins

import (
	"encoding/json"
	"fmt"

	"golang.design/x/clipboard"
	"gopkg.in/yaml.v3"

	"cmds/utils"
)

func FormatJson() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	doc := map[string]any{}
	err = json.Unmarshal(text, &doc)
	if err != nil {
		utils.Notify(err)
		return
	}
	formatedText, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		utils.Notify(err)
		return
	}
	utils.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	<-changed
	utils.Notify("previous clipboard expired")
}

func FormatSql() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	cmd := `
import sqlparse
print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='lower'))
`
	cmd = fmt.Sprintf(cmd, string(text))
	stdout, _, err := utils.RunScript("python", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	formatedText := stdout
	utils.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	<-changed
	utils.Notify("previous clipboard expired")
}

func FormatYaml() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	doc := map[any]any{}
	err = yaml.Unmarshal(text, &doc)
	if err != nil {
		utils.Notify(err)
		return
	}
	formatedText, err := yaml.Marshal(&doc)
	if err != nil {
		utils.Notify(err)
		return
	}
	utils.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	<-changed
	utils.Notify("previous clipboard expired")
}
