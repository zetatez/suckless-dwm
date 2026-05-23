package svc

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func (s *Service) bluetoothDisconnect(mac string) error {
	_, stderr, err := runScript("bash", fmt.Sprintf("bluetoothctl disconnect %s", mac))
	if err != nil {
		return fmt.Errorf("disconnect failed: %s", stderr)
	}
	return nil
}

func (s *Service) bluetoothConnect(mac string) error {
	for _, subcmd := range []string{
		fmt.Sprintf("bluetoothctl pair %s", mac),
		fmt.Sprintf("bluetoothctl trust %s", mac),
		fmt.Sprintf("bluetoothctl connect %s", mac),
	} {
		_, stderr, err := runScript("bash", subcmd)
		if err != nil {
			return fmt.Errorf("%s failed: %s", subcmd, stderr)
		}
	}
	return nil
}

func (s *Service) bluetoothScan() ([]string, error) {
	cmd := exec.Command("bluetoothctl")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	_, _ = stdin.Write([]byte("scan on\n"))

	scanner := bufio.NewScanner(stdout)
	found := make(map[string]string)
	re := regexp.MustCompile(`Device\s+([0-9A-F:]{17})\s+(.+)$`)

	scanTimeout := 6 * time.Second
	timer := time.After(scanTimeout)
	scanDone := make(chan struct{})
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if m := re.FindStringSubmatch(line); m != nil {
				found[m[1]] = m[2]
			}
		}
		close(scanDone)
	}()

	select {
	case <-timer:
	case <-scanDone:
	}

	_, _ = stdin.Write([]byte("scan off\nexit\n"))
	_ = cmd.Wait()

	var list []string
	for mac, name := range found {
		list = append(list, fmt.Sprintf("%s %s", mac, name))
	}
	return list, nil
}

func (s *Service) SysWifiConnect() error {
	out, _, err := runScript("bash", "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq|rofi -dmenu -p 'connect to wifi'")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	ssid := strings.TrimSpace(out)

	out, _, err = runScript("bash", "rofi -dmenu -p 'password' < /dev/null")
	if err != nil {
		return nil
	}
	password := strings.TrimSpace(out)

	return s.wifiConnect(ssid, password)
}

func (s *Service) wifiConnect(ssid, password string) error {
	cmd := fmt.Sprintf("nmcli device wifi connect '%s' password '%s'", ssid, password)
	_, stderr, err := runScript("bash", cmd)
	if err != nil {
		return fmt.Errorf("connect wifi failed: %s", stderr)
	}
	return nil
}

func (s *Service) SysBluetoothConnect() error {
	out, _, err := runScript("bash", "bluetoothctl devices | rofi -dmenu -p 'connect to'")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) < 2 {
		return nil
	}
	return s.bluetoothConnect(parts[1])
}

func (s *Service) SysBluetoothDisconnect() error {
	out, _, err := runScript("bash", "bluetoothctl info | grep 'Device ' | rofi -dmenu -p 'disconnect from'")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) < 2 {
		return nil
	}
	return s.bluetoothDisconnect(parts[1])
}

func (s *Service) SysBluetoothScanConnect() error {
	// scan
	devices, err := s.bluetoothScan()
	if err != nil {
		return fmt.Errorf("scan bluetooth: %w", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("no bluetooth devices found")
	}

	input := strings.Join(devices, "\n")
	out, _, err := runScript("bash", fmt.Sprintf("echo '%s' | rofi -dmenu -p 'connect bluetooth'", strings.ReplaceAll(input, "'", "'\"'\"'")))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) < 1 {
		return nil
	}
	mac := parts[0]

	return s.bluetoothConnect(mac)
}
