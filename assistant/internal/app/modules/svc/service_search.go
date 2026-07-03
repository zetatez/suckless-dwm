package svc

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

func (s *Service) Search() error {
	names := make([]string, 0, len(searchActions))
	for name := range searchActions {
		names = append(names, name)
	}
	sort.Strings(names)
	list := strings.Join(names, "\n")
	tmpf := path.Join(os.TempDir(), "assistant-search-actions")
	if err := os.WriteFile(tmpf, []byte(list), 0o644); err != nil {
		return fmt.Errorf("write action list: %w", err)
	}
	out, _, err := runScript("bash", fmt.Sprintf("rofi -dmenu -p 'search' < %s", tmpf))
	_ = os.Remove(tmpf)
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	return s.runAction(strings.TrimSpace(out))
}

func (s *Service) runAction(action string) error {
	if fn, ok := searchActions[action]; ok {
		return fn(s)
	}
	url := "https://www.google.com/search?q=" + strings.ReplaceAll(action, " ", "+")
	return s.OpenURL("chrome", url)
}

var searchActions = map[string]func(*Service) error{
	"handle clipboard":               func(s *Service) error { _, err := s.HandleClipboard(); return err },
	"format json":                    func(s *Service) error { _, err := s.Format("json"); return err },
	"format sql":                     func(s *Service) error { _, err := s.Format("sql"); return err },
	"format yaml":                    func(s *Service) error { _, err := s.Format("yaml"); return err },
	"format go":                      func(s *Service) error { _, err := s.Format("go"); return err },
	"get cur datetime":               func(s *Service) error { s.GetDatetime(); return nil },
	"get cur unix sec":               func(s *Service) error { s.GetUnixSec(); return nil },
	"get ip address":                 func(s *Service) error { _, err := s.GetIP(); return err },
	"send clipboard to feishu robot": func(s *Service) error { return s.SendToFeishu() },
	"screenshot":                     func(s *Service) error { _, err := s.Screenshot(); return err },
	"ssh to":                         func(s *Service) error { return s.SysSSHConnect() },
	"sys bluetooth connect":          func(s *Service) error { return s.SysBluetoothConnect() },
	"sys bluetooth disconnect":       func(s *Service) error { return s.SysBluetoothDisconnect() },
	"sys bluetooth scan and connect": func(s *Service) error { return s.SysBluetoothScanConnect() },
	"sys display":                    func(s *Service) error { return s.SysDisplay() },
	"sys shortcuts":                  func(s *Service) error { return s.SysShortcut() },
	"sys toggle keyboard light":      func(s *Service) error { _, err := s.SysKeyboardLight(); return err },
	"sys wifi connect":               func(s *Service) error { return s.SysWifiConnect() },
	"convert datetime to unix sec":   func(s *Service) error { _, err := s.ConvertDatetime("datetime", "unix"); return err },
	"convert unix sec to datetime":   func(s *Service) error { _, err := s.ConvertDatetime("unix", "datetime"); return err },
	"snip fzf":                       func(s *Service) error { return s.SnipFzf() },
	"snip create":                    func(s *Service) error { return s.SnipCreate("") },
	"translate clipboard":            func(s *Service) error { _, err := s.TranslateClipboard(); return err },
	"search books online": func(s *Service) error {
		q, err := s.rofiPrompt("search books")
		if err != nil || q == "" {
			return err
		}
		return s.SearchBooksOnline(q)
	},
	"search videos online": func(s *Service) error {
		q, err := s.rofiPrompt("search videos")
		if err != nil || q == "" {
			return err
		}
		return s.SearchVideosOnline(q)
	},
	"search from web": func(s *Service) error {
		q, err := s.rofiPrompt("search web")
		if err != nil || q == "" {
			return err
		}
		return s.SearchWeb(q)
	},
}
