package svc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"assistant/internal/bootstrap/psl"
	"assistant/pkg/utils"
)

func (s *Service) GitLogShow(dir string) error {
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".config")
	}

	repo, err := utils.OpenGitRepo(context.Background(), dir)
	if err != nil {
		return fmt.Errorf("open git repo: %w", err)
	}

	term := psl.GetConfig().Svc.DefaultTerminal
	script := fmt.Sprintf(`cd '%s' && git log --pretty=oneline | fzf --prompt='git show>' --ansi --preview 'git show --color=always {1}' --select-1 --exit-0 | awk '{print $1}' | xargs -o git show`, repo.Root)
	escaped := strings.ReplaceAll(script, "'", "'\\''")
	cmd := fmt.Sprintf("%s -e bash -c '%s'", term, escaped)
	if err := startScript("bash", cmd); err != nil {
		return fmt.Errorf("launch git show: %w", err)
	}
	return nil
}
