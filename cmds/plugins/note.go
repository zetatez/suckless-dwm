package plugins

import (
	"fmt"
	"os"
	"path"
	"time"

	"cmds/sugar"
)

func NoteDiary() {
	dateStr := time.Now().Format(time.DateOnly)
	fileDir := path.Join(os.Getenv("HOME"), "obsidian", "diary")
	filePath := path.Join(fileDir, dateStr+".md")
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "\n### %s\n\n", dateStr)
		f.Close()
	}
	sugar.Toggle(fmt.Sprintf("st -e nvim +':norm G' '%s'", filePath))
}

func NoteTimeline() {
	t := time.Now()
	dateStr := t.Format(time.DateOnly)
	datetimeStr := t.Format(time.DateTime)
	fileDir := path.Join(os.Getenv("HOME"), "obsidian", "timeline")
	filePath := path.Join(fileDir, dateStr+".md")
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "\n## %s\n\n", dateStr)
		f.Close()
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0o755)
	if err != nil {
		sugar.Notify(err)
		return
	}
	fmt.Fprintf(f, "\n### %s\n\n", datetimeStr)
	f.Close()
	sugar.Toggle(
		fmt.Sprintf("st -e nvim +':norm G' '%s'", filePath),
	)
}

func NoteFlashCard() {
	t := time.Now()
	fileDir := path.Join(os.Getenv("HOME"), "obsidian", "flash-card")
	filePath := path.Join(
		fileDir,
		t.Format("2006-01-02.15.04.05.000000000")+".md",
	)
	if !sugar.IsDirExists(fileDir) {
		if err := os.Mkdir(fileDir, 0o755); err != nil {
			sugar.Notify(err)
			return
		}
	}
	if !sugar.IsFileExists(filePath) {
		f, err := os.Create(filePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		fmt.Fprintf(f, "### %s\n\n", t.Format(time.DateTime))
		f.Close()
	}
	_, _, err := sugar.NewExecService().RunScriptShell(fmt.Sprintf("st -e nvim +':norm G' '%s'", filePath))
	if err != nil {
		sugar.Notify(err)
	}
}
