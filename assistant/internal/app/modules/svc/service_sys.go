package svc

import (
	"assistant/pkg/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"assistant/internal/bootstrap/psl"
	"assistant/pkg/dwmblocknotify"
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
	out, _, err := utils.RunScript("bash", fmt.Sprintf("printf '%%s\n' %s | rofi -dmenu -p 'power'", strings.Join(list, " ")))
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	action := strings.TrimSpace(out)
	cmd, ok := shortcuts[action]
	if !ok {
		return fmt.Errorf("unknown action: %s", action)
	}
	_, _, err = utils.RunScript("bash", cmd)
	return err
}

func (s *Service) SysDisplay() error {
	cfg := psl.GetConfig().Settings
	primaryMonitor := cfg.DefaultMonitor
	stdout, _, err := utils.RunScript("bash", fmt.Sprintf("xrandr|grep ' connected'|grep -v '%s'|awk '{print $1}'", primaryMonitor))
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
	out, _, err := utils.RunScript("bash", fmt.Sprintf(`printf '%%s\n' %s | rofi -dmenu -p 'screen strategy'`, strings.Join(list, " ")))
	if err != nil {
		dwmblocknotify.PUT(fmt.Sprintf("display: rofi failed: %v", err), 3*time.Second)
		s.runDisplayFallback(displayCmds["default"], cfg.DirWallpaper)
		return nil
	}
	layout := strings.TrimSpace(out)
	if layout == "" {
		return nil
	}

	cmd, ok := displayCmds[layout]
	if !ok {
		dwmblocknotify.PUT(fmt.Sprintf("display: unknown layout %q, falling back", layout), 3*time.Second)
		s.runDisplayFallback(displayCmds["default"], cfg.DirWallpaper)
		return nil
	}

	if _, stderr, err := utils.RunScript("bash", cmd); err != nil {
		return fmt.Errorf("set display: %s", stderr)
	}

	if _, _, err := utils.RunScript("bash", fmt.Sprintf("feh --bg-fill %s", cfg.DirWallpaper)); err != nil {
		dwmblocknotify.PUT(fmt.Sprintf("display: feh wallpaper failed: %v", err), 3*time.Second)
	}
	return nil
}

// runDisplayFallback executes a default xrandr/feh pair and surfaces any
// failure via notify. Used when the user cancels rofi or picks an unknown
// layout — we still want a working display, but the operator should know if
// the fallback also failed.
func (s *Service) runDisplayFallback(xrandrCmd, wallpaperDir string) {
	if _, _, err := utils.RunScript("bash", xrandrCmd); err != nil {
		dwmblocknotify.PUT(fmt.Sprintf("display fallback: xrandr failed: %v", err), 3*time.Second)
	}
	if _, _, err := utils.RunScript("bash", fmt.Sprintf("feh --bg-fill %s", wallpaperDir)); err != nil {
		dwmblocknotify.PUT(fmt.Sprintf("display fallback: feh failed: %v", err), 3*time.Second)
	}
}

func (s *Service) SysKeyboardLight() (string, error) {
	kbdPath := psl.GetConfig().Settings.PathKeyboardBrightness
	data, err := os.ReadFile(kbdPath)
	if err != nil {
		return "", fmt.Errorf("read keyboard brightness: %w", err)
	}
	current := strings.TrimSpace(string(data))
	newVal := "1"
	if current == "1" {
		newVal = "0"
	}
	_, _, err = utils.RunScript("bash", fmt.Sprintf("sudo sh -c 'echo %s > %s'", newVal, kbdPath))
	if err != nil {
		return "", fmt.Errorf("set keyboard brightness: %w", err)
	}
	return newVal, nil
}

func (s *Service) SysVolumeUp() error {
	_, _, err := utils.RunScript("bash", "amixer set Master unmute && amixer set Master 5%+")
	if err == nil {
		dwmblocknotify.PUT("vol +5%", 1*time.Second)
	}
	return err
}

func (s *Service) SysVolumeDown() error {
	_, _, err := utils.RunScript("bash", "amixer set Master unmute && amixer set Master 5%-")
	if err == nil {
		dwmblocknotify.PUT("vol -5%", 1*time.Second)
	}
	return err
}

func (s *Service) SysVolumeToggle() error {
	_, _, err := utils.RunScript("bash", "amixer set Master toggle")
	if err == nil {
		dwmblocknotify.PUT("vol toggle", 1*time.Second)
	}
	return err
}

func (s *Service) SysMicroUp() error {
	_, _, err := utils.RunScript("bash", "amixer set Capture 5%+")
	if err == nil {
		dwmblocknotify.PUT("micro +5%", 1*time.Second)
	}
	return err
}

func (s *Service) SysMicroDown() error {
	_, _, err := utils.RunScript("bash", "amixer set Capture 5%-")
	if err == nil {
		dwmblocknotify.PUT("micro -5%", 1*time.Second)
	}
	return err
}

func (s *Service) SysMicroToggle() error {
	_, _, err := utils.RunScript("bash", "amixer set Capture toggle")
	if err == nil {
		dwmblocknotify.PUT("micro toggle", 1*time.Second)
	}
	return err
}

func (s *Service) SysDisplayLightUp() error {
	_, _, err := utils.RunScript("bash", "sudo light -A 1")
	if err == nil {
		dwmblocknotify.PUT("light +1%", 1*time.Second)
	}
	return err
}

func (s *Service) SysDisplayLightDown() error {
	_, _, err := utils.RunScript("bash", "sudo light -N 1 && sudo light -U 1")
	if err == nil {
		dwmblocknotify.PUT("light -1%", 1*time.Second)
	}
	return err
}

func (s *Service) SysReset() error {
	dwmblocknotify.PUT("reset sys default", 1*time.Second)

	out, _, _ := utils.RunScript("bash", "amixer scontrols")
	hasCtl := func(name string) bool {
		return strings.Contains(out, fmt.Sprintf("'%s'", name))
	}

	type step struct{ cmd, msg string }
	steps := []step{{"sudo light -S 48", "set light to 48%"}}
	if hasCtl("Master") {
		steps = append(steps,
			step{"amixer set Master unmute", "amixer set Master unmute"},
			step{"amixer set Master 80%", "set master volume to 80%"},
		)
	}
	if hasCtl("Speaker") {
		steps = append(steps,
			step{"amixer set Speaker unmute", "amixer set Speaker unmute"},
			step{"amixer set Speaker 64%", "set speaker volume to 64%"},
		)
	}
	if hasCtl("Capture") {
		steps = append(steps,
			step{"amixer set Capture cap", "amixer set Capture cap"},
			step{"amixer set Capture 64%", "set capture volume to 64%"},
		)
	}
	if hasCtl("Headphone") {
		steps = append(steps, step{"amixer set Headphone 64%", "set headphone volume to 64%"})
	}
	steps = append(steps, step{"xset r rate 158 128", "set keyboard rate to 158 128"})

	for _, st := range steps {
		if _, _, err := utils.RunScript("bash", st.cmd); err != nil {
			dwmblocknotify.PUT(fmt.Sprintf("%s failed: %v", st.cmd, err), 3*time.Second)
		} else {
			dwmblocknotify.PUT(st.msg, 1*time.Second)
		}
	}

	return nil
}

func (s *Service) SysKill() error {
	term := psl.GetConfig().Settings.DefaultTerminal
	tmpl := `%s -e zsh -c 'ps -ef | fzf --prompt="kill -9 >" --select-1 --exit-0 | awk "{print \$2}" | xargs -r kill -9'`
	return utils.StartScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) SysWifiConnect() error {
	out, _, err := utils.RunScript("bash", "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq|rofi -dmenu -p 'connect to wifi'")
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	ssid := strings.TrimSpace(out)

	out, _, err = utils.RunScript("bash", "rofi -dmenu -p 'password' < /dev/null")
	if err != nil {
		return nil
	}
	password := strings.TrimSpace(out)

	return s.wifiConnect(ssid, password)
}

func (s *Service) wifiConnect(ssid, password string) error {
	cmd := fmt.Sprintf("nmcli device wifi connect '%s' password '%s'", ssid, password)
	_, stderr, err := utils.RunScript("bash", cmd)
	if err != nil {
		return fmt.Errorf("connect wifi failed: %s", stderr)
	}
	return nil
}

func (s *Service) SysBluetoothConnect() error {
	out, _, err := utils.RunScript("bash", "bluetoothctl devices | rofi -dmenu -p 'connect to'")
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
	out, _, err := utils.RunScript("bash", "bluetoothctl info | grep 'Device ' | rofi -dmenu -p 'disconnect from'")
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
	out, _, err := utils.RunScript("bash", fmt.Sprintf("echo '%s' | rofi -dmenu -p 'connect bluetooth'", strings.ReplaceAll(input, "'", "'\"'\"'")))
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
	_, stderr, err := utils.RunScript("bash", fmt.Sprintf("bluetoothctl disconnect %s", mac))
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
		_, stderr, err := utils.RunScript("bash", subcmd)
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
