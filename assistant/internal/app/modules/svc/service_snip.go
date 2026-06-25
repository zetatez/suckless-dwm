package svc

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) SnipFzf() error {
	snipDir := psl.GetConfig().Svc.SnipDir
	if _, err := os.Stat(snipDir); err != nil {
		return fmt.Errorf("snippet dir not found: %s", snipDir)
	}

	tmpf := path.Join(os.TempDir(), "snip-fzf-selected")
	_ = os.Remove(tmpf)
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `
cd %s && selected=$(
  find . -type f | sed "s|^\./||" |
  fzf --prompt="Snip> " --height=100%% --border \
      --preview="bat --style=plain --color=always {} 2>/dev/null || cat {}" \
      --preview-window=right:60%%
) && [ -n "$selected" ] && printf "%%s" "$selected" > %s
`
	script := fmt.Sprintf(tmpl, snipDir, tmpf)
	cmd := fmt.Sprintf("%s -e sh -c '%s'", term, script)
	if _, _, err := runScript("bash", cmd); err != nil {
		return fmt.Errorf("launch fzf: %w", err)
	}

	<-time.After(300 * time.Millisecond)
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

	if _, err := s.pushClipboard(string(content), fmt.Sprintf("Snip copied: %s", file)); err != nil {
		return err
	}
	return nil
}

func (s *Service) SnipCreate(name string) error {
	snipDir := psl.GetConfig().Svc.SnipDir
	if err := os.MkdirAll(snipDir, 0o755); err != nil {
		return fmt.Errorf("create snippet dir: %w", err)
	}

	if name == "" {
		out, _, err := runScript("bash", "rofi -dmenu -p 'Snippet name' < /dev/null")
		if err != nil || strings.TrimSpace(out) == "" {
			return nil
		}
		name = strings.TrimSpace(out)
	}

	filePath := path.Join(snipDir, name)
	if err := os.MkdirAll(path.Dir(filePath), 0o755); err != nil {
		return fmt.Errorf("create snippet subdir: %w", err)
	}
	term := psl.GetConfig().Svc.DefaultTerminal
	if err := startScript("bash", fmt.Sprintf("%s -e nvim '%s'", term, filePath)); err != nil {
		return fmt.Errorf("launch nvim: %w", err)
	}
	s.notify(fmt.Sprintf("Snip created: %s", name))
	return nil
}
