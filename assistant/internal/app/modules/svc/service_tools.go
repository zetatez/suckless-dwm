package svc

import (
	"encoding/json"
	"fmt"
	"go/format"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"

	"gopkg.in/yaml.v3"
)

func (s *Service) Format(language string) (string, error) {
	content, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	if content == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	var result string

	switch language {
	case "json":
		var doc any
		if err := json.Unmarshal([]byte(content), &doc); err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
		formatted, e := json.MarshalIndent(doc, "", "  ")
		if e != nil {
			return "", fmt.Errorf("format JSON failed: %w", e)
		}
		result = string(formatted)
	case "yaml":
		var doc any
		if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
			return "", fmt.Errorf("invalid YAML: %w", err)
		}
		formatted, e := yaml.Marshal(&doc)
		if e != nil {
			return "", fmt.Errorf("format YAML failed: %w", e)
		}
		result = string(formatted)
	case "sql":
		stdout, stderr, e := runScript("python", fmt.Sprintf(
			`import sqlparse; print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='lower'))`, content))
		if e != nil {
			return "", fmt.Errorf("format SQL failed: %s", stderr)
		}
		result = strings.TrimSpace(stdout)
	case "go":
		formatted, e := format.Source([]byte(content))
		if e != nil {
			return "", fmt.Errorf("format Go failed: %w", e)
		}
		result = string(formatted)
	default:
		return "", fmt.Errorf("unsupported language: %s, available: json, yaml, sql, go", language)
	}

	s.writeClipboard(result)
	s.notify(fmt.Sprintf("format %s success", language))
	go func() {
		time.Sleep(30 * time.Second)
		s.notify("previous clipboard expired")
	}()

	return result, nil
}

func (s *Service) Note(noteType string) error {
	fileDir := expandHome(psl.GetConfig().Svc.WorkingLogbookDir)
	if err := os.MkdirAll(fileDir, 0o755); err != nil {
		return fmt.Errorf("create logbook dir: %w", err)
	}

	var filePath string
	var header string
	now := time.Now()

	switch noteType {
	case "todo":
		filePath = path.Join(fileDir, "TODO.md")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if f, err := os.Create(filePath); err == nil {
				_, _ = fmt.Fprintf(f, "\n## ToDo\n\n")
				_ = f.Close()
			}
		}
		header = fmt.Sprintf("\n- [ ] %s", now.Format(time.DateTime))
	case "scripts":
		filePath = path.Join(fileDir, "scripts.md")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if f, err := os.Create(filePath); err == nil {
				_, _ = fmt.Fprintf(f, "\n## Scripts\n\n")
				_ = f.Close()
			}
		}
		header = fmt.Sprintf("\n\n### %s", now.Format(time.DateTime))
	case "monthly-work":
		dateStr := now.Format("2006-01")
		filePath = path.Join(fileDir, dateStr+".md")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if f, err := os.Create(filePath); err == nil {
				_, _ = fmt.Fprintf(f, "\n## %s\n\n", dateStr)
				_ = f.Close()
			}
		}
		header = fmt.Sprintf("\n### %s\n\n", now.Format(time.DateTime))
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
	_ = f.Close()

	term := psl.GetConfig().Svc.DefaultTerminal
	_, _, _ = runScript("bash", fmt.Sprintf("%s -e nvim +$ '%s'", term, filePath))
	return nil
}

func (s *Service) GetDatetime() map[string]string {
	now := time.Now()
	result := map[string]string{
		"datetime":   now.Format(time.DateTime),
		"unix":       fmt.Sprintf("%d", now.Unix()),
		"unix_milli": fmt.Sprintf("%d", now.UnixMilli()),
		"date":       now.Format("2006-01-02"),
		"time":       now.Format("15:04:05"),
	}
	s.writeClipboard(result["datetime"])
	s.notify(fmt.Sprintf("get success: %s", result["datetime"]))
	go func() {
		time.Sleep(30 * time.Second)
		s.notify("previous clipboard expired")
	}()
	return result
}

func (s *Service) GetCurUnixSec() string {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s.writeClipboard(now)
	s.notify(fmt.Sprintf("get success: %s", now))
	go func() {
		time.Sleep(30 * time.Second)
		s.notify("previous clipboard expired")
	}()
	return now
}

func (s *Service) SearchBooksOnline(query string) error {
	q := strings.ReplaceAll(query, " ", "+")
	urls := []string{
		"https://openlibrary.org/search?q=" + q,
		"https://z-lib.id/s?q=" + q,
	}
	for _, u := range urls {
		_ = s.OpenURL("chrome", u)
	}
	return nil
}

func (s *Service) SearchVideosOnline(query string) error {
	q := strings.ReplaceAll(query, " ", "+")
	urls := []string{
		"https://search.bilibili.com/all?keyword=" + q,
		"https://www.youtube.com/results?search_query=" + q,
	}
	for _, u := range urls {
		_ = s.OpenURL("chrome", u)
	}
	return nil
}

func (s *Service) ConvertDatetime(from, to string) (string, error) {
	clip, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	input := strings.TrimSpace(clip)
	if input == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	var t time.Time

	switch from {
	case "datetime":
		t, err = time.Parse(time.DateTime, strings.TrimSpace(input))
	case "unix":
		sec, e := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if e != nil {
			return "", fmt.Errorf("invalid unix timestamp: %w", e)
		}
		t = time.Unix(sec, 0)
	case "unix_milli":
		ms, e := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if e != nil {
			return "", fmt.Errorf("invalid unix milli timestamp: %w", e)
		}
		t = time.UnixMilli(ms)
	default:
		return "", fmt.Errorf("unsupported from format: %s, available: datetime, unix, unix_milli", from)
	}
	if err != nil {
		return "", fmt.Errorf("parse input failed: %w", err)
	}

	switch to {
	case "datetime":
		result := t.Format(time.DateTime)
		s.writeClipboard(result)
		s.notify(fmt.Sprintf("transfer success: %s", result))
		go func() { time.Sleep(30 * time.Second); s.notify("previous clipboard expired") }()
		return result, nil
	case "unix":
		result := fmt.Sprintf("%d", t.Unix())
		s.writeClipboard(result)
		s.notify(fmt.Sprintf("transfer success: %s", result))
		go func() { time.Sleep(30 * time.Second); s.notify("previous clipboard expired") }()
		return result, nil
	case "unix_milli":
		result := fmt.Sprintf("%d", t.UnixMilli())
		s.writeClipboard(result)
		s.notify(fmt.Sprintf("transfer success: %s", result))
		go func() { time.Sleep(30 * time.Second); s.notify("previous clipboard expired") }()
		return result, nil
	default:
		return "", fmt.Errorf("unsupported to format: %s, available: datetime, unix, unix_milli", to)
	}
}

func (s *Service) GetIP(iface string) ([]string, error) {
	if iface == "" {
		iface = psl.GetConfig().App.Interface
	}
	netIface, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, fmt.Errorf("get interface %s: %w", iface, err)
	}
	addrs, err := netIface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("get addresses: %w", err)
	}
	var ips []string
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if !ip.IsLoopback() && ip.To4() != nil {
			ips = append(ips, ip.String())
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPv4 addresses found on %s", iface)
	}
	s.writeClipboard(ips[0])
	s.notify(fmt.Sprintf("get success: %s", ips[0]))
	go func() { time.Sleep(30 * time.Second); s.notify("previous clipboard expired") }()
	return ips, nil
}
