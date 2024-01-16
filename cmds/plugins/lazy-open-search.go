package plugins

import (
	"cmds/sugar"
)

func LazyOpenSearchFile() {
	cmd := `st -e lazy-open-search-file`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchBook() {
	cmd := `st -e lazy-open-search-book`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchWiki() {
	cmd := `st -e lazy-open-search-wiki`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchMedia() {
	cmd := `st -e lazy-open-search-media`
	sugar.NewExecService().RunScriptShell(cmd)
}

func LazyOpenSearchFileContent() {
	cmd := `st -e lazy-open-search-file-content`
	sugar.NewExecService().RunScriptShell(cmd)
}
