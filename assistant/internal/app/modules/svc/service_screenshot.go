package svc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func (s *Service) Screenshot() (string, error) {
	tool := "flameshot"
	if _, err := exec.LookPath(tool); err != nil {
		return "", fmt.Errorf("screenshot tool not found: %s", tool)
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	dir = filepath.Join(dir, "Pictures", "screenshots")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create screenshot dir: %w", err)
	}

	now := time.Now()
	filename := fmt.Sprintf("screenshot.%s.jpeg", now.Format("2006.01.02.15.04.05"))
	path := filepath.Join(dir, filename)

	cmd := exec.Command(tool, "full", "-p", path)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	if err := s.writeClipboard(path); err != nil {
		s.logger.WithError(err).Warn("write clipboard failed")
	}

	return path, nil
}
