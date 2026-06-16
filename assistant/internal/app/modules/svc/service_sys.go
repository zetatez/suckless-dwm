package svc

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) SysShortcut() error {
	shortcuts := map[string]string{
		"suspend":     "systemctl suspend",
		"poweroff":    "systemctl poweroff",
		"reboot":      "systemctl reboot",
		"off-display": "sleep .5; xset dpms force off",
		"slock":       "slock",
	}
	list := make([]string, 0, len(shortcuts))
	for k := range shortcuts {
		list = append(list, k)
	}
	out, _, err := runScript("bash", fmt.Sprintf("printf '%%s\n' %s | rofi -dmenu -p 'power'", strings.Join(list, " ")))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	action := strings.TrimSpace(out)
	cmd, ok := shortcuts[action]
	if !ok {
		return fmt.Errorf("unknown action: %s", action)
	}
	_, _, err = runScript("bash", cmd)
	return err
}

func (s *Service) SysDisplay() error {
	cfg := psl.GetConfig().Svc
	primaryMonitor := cfg.PrimaryMonitor
	stdout, _, err := runScript("bash", fmt.Sprintf("xrandr|grep ' connected'|grep -v '%s'|awk '{print $1}'", primaryMonitor))
	if err != nil {
		return fmt.Errorf("detect monitors: %w", err)
	}
	secondMonitor := strings.TrimSpace(stdout)
	if secondMonitor == "" {
		return fmt.Errorf("only one monitor detected")
	}

	displayCmds := map[string]string{
		"default":               fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"clone":                 fmt.Sprintf("xrandr --output %s --mode 1920x1080", secondMonitor),
		"primary_only":          fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"second_only":           fmt.Sprintf("xrandr --output %s --auto --output %s --off", secondMonitor, primaryMonitor),
		"left_of":               fmt.Sprintf("xrandr --output %s --auto --left-of %s --auto", secondMonitor, primaryMonitor),
		"right_of":              fmt.Sprintf("xrandr --output %s --auto --right-of %s --auto", secondMonitor, primaryMonitor),
		"above":                 fmt.Sprintf("xrandr --output %s --auto --above %s --auto", secondMonitor, primaryMonitor),
		"below":                 fmt.Sprintf("xrandr --output %s --auto --below %s --auto", secondMonitor, primaryMonitor),
		"rotate_left_left_of":   fmt.Sprintf("xrandr --output %s --auto --rotate left --left-of %s --auto", secondMonitor, primaryMonitor),
		"rotate_right_left_of":  fmt.Sprintf("xrandr --output %s --auto --rotate right --left-of %s --auto", secondMonitor, primaryMonitor),
		"rotate_left_right_of":  fmt.Sprintf("xrandr --output %s --auto --rotate left --right-of %s --auto", secondMonitor, primaryMonitor),
		"rotate_right_right_of": fmt.Sprintf("xrandr --output %s --auto --rotate right --right-of %s --auto", secondMonitor, primaryMonitor),
	}

	list := make([]string, 0, len(displayCmds))
	for k := range displayCmds {
		list = append(list, k)
	}
	out, _, err := runScript("bash", fmt.Sprintf(`printf '%%s\n' %s | rofi -dmenu -p 'screen strategy'`, strings.Join(list, " ")))
	if err != nil {
		s.notify(fmt.Sprintf("display: rofi failed: %v", err))
		s.runDisplayFallback(displayCmds["default"], cfg.DirWallpaper)
		return nil
	}
	layout := strings.TrimSpace(out)
	if layout == "" {
		return nil
	}

	cmd, ok := displayCmds[layout]
	if !ok {
		s.notify(fmt.Sprintf("display: unknown layout %q, falling back", layout))
		s.runDisplayFallback(displayCmds["default"], cfg.DirWallpaper)
		return nil
	}

	if _, stderr, err := runScript("bash", cmd); err != nil {
		return fmt.Errorf("set display: %s", stderr)
	}

	if _, _, err := runScript("bash", fmt.Sprintf("feh --bg-fill %s", cfg.DirWallpaper)); err != nil {
		s.notify(fmt.Sprintf("display: feh wallpaper failed: %v", err))
	}
	return nil
}

// runDisplayFallback executes a default xrandr/feh pair and surfaces any
// failure via notify. Used when the user cancels rofi or picks an unknown
// layout — we still want a working display, but the operator should know if
// the fallback also failed.
func (s *Service) runDisplayFallback(xrandrCmd, wallpaperDir string) {
	if _, _, err := runScript("bash", xrandrCmd); err != nil {
		s.notify(fmt.Sprintf("display fallback: xrandr failed: %v", err))
	}
	if _, _, err := runScript("bash", fmt.Sprintf("feh --bg-fill %s", wallpaperDir)); err != nil {
		s.notify(fmt.Sprintf("display fallback: feh failed: %v", err))
	}
}

func (s *Service) SysKeyboardLight() (string, error) {
	kbdPath := psl.GetConfig().Svc.KeyboardBrightnessPath
	data, err := os.ReadFile(kbdPath)
	if err != nil {
		return "", fmt.Errorf("read keyboard brightness: %w", err)
	}
	current := strings.TrimSpace(string(data))
	newVal := "1"
	if current == "1" {
		newVal = "0"
	}
	_, _, err = runScript("bash", fmt.Sprintf("sudo sh -c 'echo %s > %s'", newVal, kbdPath))
	if err != nil {
		return "", fmt.Errorf("set keyboard brightness: %w", err)
	}
	return newVal, nil
}

func (s *Service) SysVolumeUp() error {
	_, _, err := runScript("bash", "amixer set Master unmute && amixer set Master 5%+")
	if err == nil {
		s.notify("vol +5%")
	}
	return err
}

func (s *Service) SysVolumeDown() error {
	_, _, err := runScript("bash", "amixer set Master unmute && amixer set Master 5%-")
	if err == nil {
		s.notify("vol -5%")
	}
	return err
}

func (s *Service) SysVolumeToggle() error {
	_, _, err := runScript("bash", "amixer set Master toggle")
	if err == nil {
		s.notify("vol toggle")
	}
	return err
}

func (s *Service) SysMicroUp() error {
	_, _, err := runScript("bash", "amixer set Capture 5%+")
	if err == nil {
		s.notify("micro +5%")
	}
	return err
}

func (s *Service) SysMicroDown() error {
	_, _, err := runScript("bash", "amixer set Capture 5%-")
	if err == nil {
		s.notify("micro -5%")
	}
	return err
}

func (s *Service) SysMicroToggle() error {
	_, _, err := runScript("bash", "amixer set Capture toggle")
	if err == nil {
		s.notify("micro toggle")
	}
	return err
}

func (s *Service) SysDisplayLightUp() error {
	_, _, err := runScript("bash", "sudo light -A 1")
	if err == nil {
		s.notify("light +1%")
	}
	return err
}

func (s *Service) SysDisplayLightDown() error {
	_, _, err := runScript("bash", "sudo light -N 1 && sudo light -U 1")
	if err == nil {
		s.notify("light -1%")
	}
	return err
}

func (s *Service) SysReset() error {
	s.notify("reset sys default")
	steps := []struct {
		cmd string
		msg string
	}{
		{"sudo light -S 48", "set light to 48%"},
		{"amixer set Master unmute", "amixer set Master unmute"},
		{"amixer set Speaker unmute", "amixer set Speaker unmute"},
		{"amixer set Capture cap", "amixer set Capture cap"},
		{"amixer set Master 80%", "set master volume to 80%"},
		{"amixer set Capture 64%", "set capture volume to 64%"},
		{"amixer set Speaker 64%", "set speaker volume to 64%"},
		{"amixer set Headphone 64%", "set headphone volume to 64%"},
		{"xset r rate 158 128", "set keyboard rate to 158 128"},
	}

	for _, step := range steps {
		if _, _, err := runScript("bash", step.cmd); err != nil {
			s.notify(fmt.Sprintf("%s failed: %v", step.cmd, err))
		} else {
			s.notify(step.msg)
		}
	}

	return nil
}

func (s *Service) SysKill() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e zsh -c 'ps -ef | fzf --prompt="kill -9 >" --select-1 --exit-0 | awk "{print \$2}" | xargs -r kill -9'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
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
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("bluetooth stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("bluetooth stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start bluetoothctl: %w", err)
	}
	if _, err := stdin.Write([]byte("scan on\n")); err != nil {
		return nil, fmt.Errorf("write scan on: %w", err)
	}

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

	if _, err := stdin.Write([]byte("scan off\nexit\n")); err != nil {
		s.logger.WithError(err).Warn("bluetooth write scan off failed")
	}
	_ = cmd.Wait()

	var list []string
	for mac, name := range found {
		list = append(list, fmt.Sprintf("%s %s", mac, name))
	}
	return list, nil
}
