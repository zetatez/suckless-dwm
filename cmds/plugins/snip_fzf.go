package plugins

import (
	"cmds/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
)

func SnipFzf() error {
	snipDir := os.ExpandEnv("$HOME/share/github/obsidian/.snippets")
	if _, err := os.Stat(snipDir); err != nil {
		return fmt.Errorf("snippet dir not found: %s", snipDir)
	}

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return err
	}
	defer readPipe.Close()

	// shell 脚本: 选择 snip
	script := fmt.Sprintf(`
cd %s || exit 1
selected=$(find . -type f | sed 's|^\./||' |
fzf \
  --prompt="Snip> " \
  --height=100%% \
  --border \
  --preview='bat --style=plain --color=always {} 2>/dev/null || cat {}' \
  --preview-window=right:60%%)
printf '%%s' "$selected" >&3
`, snipDir)

	cmd := exec.Command(utils.GetOSDefaultTerminal(), "-e", "sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{writePipe}

	if err := cmd.Start(); err != nil {
		writePipe.Close()
		return err
	}
	writePipe.Close()

	if err := cmd.Wait(); err != nil {
		return nil // 用户 ESC 退出时，fzf 返回非 0，直接忽略
	}

	data, err := io.ReadAll(readPipe)
	if err != nil {
		return err
	}

	file := strings.TrimSpace(string(data))
	if file == "" {
		return nil
	}

	content, err := os.ReadFile(filepath.Join(snipDir, file))
	if err != nil {
		return err
	}

	if err := clipboard.Init(); err != nil {
		return err
	}
	utils.Notify(fmt.Sprintf("Snip copied:\n%s", file))
	changed := clipboard.Write(clipboard.FmtText, content)
	<-changed
	utils.Notify("previous clipboard expired")
	return nil
}
