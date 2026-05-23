package svc

import (
	"fmt"
	"strings"
)

var searchActions = map[string]func(*Service) error{
	"handle clipboard":               func(s *Service) error { _, err := s.HandleClipboard(); return err },
	"format json":                    func(s *Service) error { _, err := s.Format("json"); return err },
	"format sql":                     func(s *Service) error { _, err := s.Format("sql"); return err },
	"format yaml":                    func(s *Service) error { _, err := s.Format("yaml"); return err },
	"format go":                      func(s *Service) error { _, err := s.Format("go"); return err },
	"get cur datetime":               func(s *Service) error { s.GetDatetime(); return nil },
	"get cur unix sec":               func(s *Service) error { s.GetCurUnixSec(); return nil },
	"get ip address":                 func(s *Service) error { _, err := s.GetIP(""); return err },
	"send clipboard to feishu robot": func(s *Service) error { return s.FeishuSend() },
	"note script":                    func(s *Service) error { return s.Note("scripts") },
	"note todo":                      func(s *Service) error { return s.Note("todo") },
	"note monthly work":              func(s *Service) error { return s.Note("monthly_work") },
	"ssh to":                         func(s *Service) error { return s.SysSSHConnect() },
	"search books online": func(s *Service) error {
		q, _, _ := runScript("bash", `printf '' | rofi -dmenu -p 'search books'`)
		if strings.TrimSpace(q) == "" {
			return nil
		}
		return s.SearchBooksOnline(strings.TrimSpace(q))
	},
	"search videos online": func(s *Service) error {
		q, _, _ := runScript("bash", `printf '' | rofi -dmenu -p 'search videos'`)
		if strings.TrimSpace(q) == "" {
			return nil
		}
		return s.SearchVideosOnline(strings.TrimSpace(q))
	},
	"search from web": func(s *Service) error {
		q, _, _ := runScript("bash", `printf '' | rofi -dmenu -p 'search web'`)
		if strings.TrimSpace(q) == "" {
			return nil
		}
		return s.SearchWeb(strings.TrimSpace(q))
	},
	"sys bluetooth connect":           func(s *Service) error { return s.SysBluetoothConnect() },
	"sys bluetooth disconnect":        func(s *Service) error { return s.SysBluetoothDisconnect() },
	"sys bluetooth scan and connect":  func(s *Service) error { return s.SysBluetoothScanConnect() },
	"sys display":                     func(s *Service) error { return s.SysDisplay() },
	"sys shortcuts":                   func(s *Service) error { return s.SysShortcut() },
	"sys toggle keyboard light":       func(s *Service) error { _, err := s.SysKeyboardLight(); return err },
	"sys wifi connect":                func(s *Service) error { return s.SysWifiConnect() },
	"conversion datetime to unix sec": func(s *Service) error { _, err := s.ConvertDatetime("datetime", "unix"); return err },
	"conversion unix sec to datetime": func(s *Service) error { _, err := s.ConvertDatetime("unix", "datetime"); return err },
	"launch inkscape":                 func(s *Service) error { return s.Launch("inkscape") },
	"launch krita":                    func(s *Service) error { return s.Launch("krita") },
	"launch obsidian":                 func(s *Service) error { return s.Launch("obsidian") },
	"launch sublime":                  func(s *Service) error { return s.Launch("subl") },
	"launch xournal":                  func(s *Service) error { return s.Launch("xournalpp") },
	"snip fzf":                        func(s *Service) error { return s.SnipFzf() },
	"snip create":                     func(s *Service) error { return s.SnipCreate("") },
	"file search":                     func(s *Service) error { return s.FileSearch() },
	"file search content":             func(s *Service) error { return s.FileSearchContent() },
	"file search book":                func(s *Service) error { return s.FileSearchBook() },
	"file search media":               func(s *Service) error { return s.FileSearchMedia() },
	"file search wiki":                func(s *Service) error { return s.FileSearchWiki() },
}

func (s *Service) runAction(action string) error {
	if fn, ok := searchActions[action]; ok {
		return fn(s)
	}
	if strings.HasPrefix(action, "toggle ") {
		return fmt.Errorf("use toggle.sh with process name: %s", strings.TrimPrefix(action, "toggle "))
	}
	return fmt.Errorf("unknown action: %s", action)
}
