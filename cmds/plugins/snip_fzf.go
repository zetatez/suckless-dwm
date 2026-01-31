package plugins

import (
	"cmds/utils"
	"fmt"
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

	// 创建临时文件保存 fzf 选择结果
	tmpFile, err := os.CreateTemp("", "snip_fzf_*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// shell 脚本: 选择 snip
	script := fmt.Sprintf(`
cd %s || exit 1
find . -type f | sed 's|^\./||' |
fzf \
  --prompt="Snip> " \
  --height=100%% \
  --border \
  --preview='bat --style=plain --color=always {} 2>/dev/null || cat {}' \
  --preview-window=right:60%% \
> %s
`, snipDir, tmpPath)

	// st 中运行 fzf
	cmd := exec.Command("st", "-e", "sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// 用户 ESC 退出时，fzf 返回非 0，这里直接忽略
		return nil
	}

	// 读取 fzf 结果
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return err
	}

	file := strings.TrimSpace(string(data))
	if file == "" {
		return nil
	}

	// 读取 snippet 内容
	content, err := os.ReadFile(filepath.Join(snipDir, file))
	if err != nil {
		return err
	}

	// 写入剪贴板
	if err := clipboard.Init(); err != nil {
		return err
	}
	utils.Notify(fmt.Sprintf("Snip copied:\n%s", file))
	changed := clipboard.Write(clipboard.FmtText, content)
	<-changed
	utils.Notify("previous clipboard expired")
	return nil
}
