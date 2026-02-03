package plugins

import (
	"encoding/json"
	"fmt"
	"go/format"
	"time"

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
	clipboard.Write(clipboard.FmtText, formatedText)
	time.Sleep(30 * time.Second)
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
	clipboard.Write(clipboard.FmtText, []byte(formatedText))
	time.Sleep(30 * time.Second)
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
	clipboard.Write(clipboard.FmtText, formatedText)
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}

func FormatGo() {
	err := clipboard.Init()
	if err != nil {
		utils.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	if len(text) == 0 {
		utils.Notify("clipboard empty or not text")
		return
	}
	formatted, err := format.Source(text)
	if err != nil {
		utils.Notify(fmt.Errorf("format failed: %w", err))
		return
	}
	utils.Notify(fmt.Sprintf("format success:\n%s", formatted))
	clipboard.Write(clipboard.FmtText, formatted)
	time.Sleep(30 * time.Second)
	utils.Notify("previous clipboard expired")
}
