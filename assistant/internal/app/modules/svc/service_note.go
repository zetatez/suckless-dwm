package svc

import (
	"assistant/pkg/utils"
	"fmt"
	"os"
	"path"
	"time"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) Note(noteType string) error {
	fileDir := psl.GetConfig().Svc.WorkingLogbookDir
	if err := os.MkdirAll(fileDir, 0o755); err != nil {
		return fmt.Errorf("create logbook dir: %w", err)
	}

	var filePath string
	var header string
	now := time.Now()

	switch noteType {
	case "todo":
		filePath = path.Join(fileDir, "TODO.md")
		if err := ensureNoteFile(filePath, "\n## ToDo\n\n"); err != nil {
			return err
		}
		header = fmt.Sprintf("\n- [ ] %s:", now.Format(time.DateTime))
	case "scripts":
		filePath = path.Join(fileDir, "scripts.md")
		if err := ensureNoteFile(filePath, "\n## Scripts\n\n"); err != nil {
			return err
		}
		header = fmt.Sprintf("\n\n### %s:", now.Format(time.DateTime))
	case "monthly-work":
		dateStr := now.Format("2006-01")
		filePath = path.Join(fileDir, dateStr+".md")
		if err := ensureNoteFile(filePath, fmt.Sprintf("\n## %s\n\n", dateStr)); err != nil {
			return err
		}
		header = fmt.Sprintf("\n\n### %s:", now.Format(time.DateTime))
	default:
		return fmt.Errorf("unknown note type: %s, available: todo, scripts, monthly-work", noteType)
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open note file: %w", err)
	}
	if _, err := f.WriteString(header); err != nil {
		_ = f.Close()
		return fmt.Errorf("write note: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close note file: %w", err)
	}

	term := psl.GetConfig().Svc.DefaultTerminal
	if _, _, err := utils.RunScript("bash", fmt.Sprintf(`%s -e nvim "+normal G$" "%s"`, term, filePath)); err != nil {
		return fmt.Errorf("launch nvim: %w", err)
	}
	return nil
}

// ensureNoteFile creates the file with the given header if it does not yet
// exist. Existing files are left untouched.
func ensureNoteFile(filePath, header string) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat note file: %w", err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create note file: %w", err)
	}
	if _, err := fmt.Fprint(f, header); err != nil {
		_ = f.Close()
		return fmt.Errorf("write note header: %w", err)
	}
	return f.Close()
}
