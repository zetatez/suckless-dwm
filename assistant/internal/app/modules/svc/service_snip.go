package svc

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) SnipFzf() error {
	homeDir, _ := os.UserHomeDir()
	snipDir := path.Join(homeDir, "git/obsidian/.snippets")
	if _, err := os.Stat(snipDir); err != nil {
		return fmt.Errorf("snippet dir not found: %s", snipDir)
	}

	tmpf := path.Join(os.TempDir(), "snip-fzf-selected")
	os.Remove(tmpf)
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `
cd %s && selected=$(
  find . -type f | sed "s|^\./||" |
  fzf --prompt="Snip> " --height=100%% --border \
      --preview="bat --style=plain --color=always {} 2>/dev/null || cat {}" \
      --preview-window=right:60%%)
) && [ -n "$selected" ] && printf "%%s" "$selected" > %s
`
	script := fmt.Sprintf(tmpl, snipDir, tmpf)
	cmd := fmt.Sprintf("%s -e sh -c '%s'", term, script)
	_, _, _ = runScript("bash", cmd)

	time.Sleep(200 * time.Millisecond)
	data, err := os.ReadFile(tmpf)
	if err != nil || len(data) == 0 {
		return nil
	}
	file := strings.TrimSpace(string(data))
	if file == "" {
		return nil
	}

	content, err := os.ReadFile(path.Join(snipDir, file))
	if err != nil {
		return err
	}

	s.writeClipboard(string(content))
	s.notify(fmt.Sprintf("Snip copied: %s", file))
	go func() {
		time.Sleep(30 * time.Second)
		s.notify("previous clipboard expired")
	}()
	return nil
}

func (s *Service) SnipCreate(name string) error {
	homeDir, _ := os.UserHomeDir()
	snipDir := path.Join(homeDir, "git/obsidian/.snippets")
	os.MkdirAll(snipDir, 0o755)

	if name == "" {
		out, _, err := runScript("bash", "rofi -dmenu -p 'Snippet name' < /dev/null")
		if err != nil || strings.TrimSpace(out) == "" {
			return nil
		}
		name = strings.TrimSpace(out)
	}

	filePath := path.Join(snipDir, name)
	os.MkdirAll(path.Dir(filePath), 0o755)
	term := psl.GetConfig().Svc.DefaultTerminal
	_ = startScript("bash", fmt.Sprintf("%s -e nvim '%s'", term, filePath))
	s.notify(fmt.Sprintf("Snip created: %s", name))
	return nil
}

func (s *Service) Search() error {
	names := make([]string, 0, len(searchActions))
	for name := range searchActions {
		names = append(names, name)
	}
	sort.Strings(names)
	list := strings.Join(names, "\n")
	tmpf := path.Join(os.TempDir(), "assistant-search-actions")
	os.WriteFile(tmpf, []byte(list), 0o644)
	out, _, err := runScript("bash", fmt.Sprintf("rofi -dmenu -p 'search' < %s", tmpf))
	os.Remove(tmpf)
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	return s.runAction(strings.TrimSpace(out))
}
