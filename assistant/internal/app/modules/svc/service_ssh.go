package svc

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) ConnectSSH(host string, port int, user, password string) error {
	if port <= 0 {
		port = 22
	}
	if user == "" {
		user = "root"
	}
	term := psl.GetConfig().Svc.DefaultTerminal
	cmd := fmt.Sprintf("%s -e sshpass -p '%s' ssh -o 'StrictHostKeyChecking no' -p %d %s@%s &", term, password, port, user, host)
	_, _, err := runScript("bash", cmd)
	return err
}

func (s *Service) SysSSHConnect() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sshConfigPath := psl.GetConfig().Svc.SSHSecretFile
	sshConfigPath = strings.ReplaceAll(sshConfigPath, "~", homeDir)

	var lines []string
	if data, err := os.ReadFile(sshConfigPath); err == nil {
		saved := map[string]struct {
			Host     string `json:"host"`
			User     string `json:"user"`
			Password string `json:"password"`
			Desc     string `json:"desc"`
		}{}
		if err := json.Unmarshal(data, &saved); err == nil {
			for _, entry := range saved {
				line := fmt.Sprintf("%s@%s  # %s", entry.User, entry.Host, entry.Desc)
				lines = append(lines, line)
			}
		}
	}

	knownHostsPath := path.Join(homeDir, ".ssh/known_hosts")
	if data, err := os.ReadFile(knownHostsPath); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 1 && fields[0] != "" {
				host := strings.TrimSpace(fields[0])
				if !strings.Contains(host, ",") {
					lines = append(lines, host)
				}
			}
		}
	}

	if len(lines) == 0 {
		return fmt.Errorf("no SSH entries found")
	}

	input := strings.Join(lines, "\n")
	out, _, err := runScript("bash", fmt.Sprintf("echo '%s' | rofi -dmenu -p 'ssh to'", strings.ReplaceAll(input, "'", "'\"'\"'")))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	choice := strings.TrimSpace(out)

	host, user, password := "", "", ""
	parts := strings.Fields(choice)
	for _, p := range parts {
		if strings.Contains(p, "@") {
			parts2 := strings.Split(p, "@")
			user = parts2[0]
			host = parts2[1]
		} else if strings.Contains(p, "#") {
			// description, skip
		} else if host == "" {
			host = p
		}
	}

	if host == "" {
		return fmt.Errorf("no host selected")
	}
	if user == "" {
		user = "root"
	}

	if data, err := os.ReadFile(sshConfigPath); err == nil {
		saved := map[string]struct {
			Host     string `json:"host"`
			User     string `json:"user"`
			Password string `json:"password"`
			Desc     string `json:"desc"`
		}{}
		if err := json.Unmarshal(data, &saved); err == nil {
			if entry, ok := saved[user+"@"+host]; ok {
				password = entry.Password
			}
		}
	}

	if password == "" {
		out2, _, err := runScript("bash", "rofi -dmenu -p 'password' < /dev/null")
		if err != nil || strings.TrimSpace(out2) == "" {
			return nil
		}
		password = strings.TrimSpace(out2)
	}

	return s.ConnectSSH(host, 22, user, password)
}
