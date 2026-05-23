package svc

import (
	"fmt"
	"os"
	"strings"

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
		_, _, _ = runScript("bash", displayCmds["default"])
		_, _, _ = runScript("bash", fmt.Sprintf("feh --bg-fill %s", expandHome(cfg.DirWallpaper)))
		return nil
	}
	layout := strings.TrimSpace(out)
	if layout == "" {
		return nil
	}

	cmd, ok := displayCmds[layout]
	if !ok {
		_, _, _ = runScript("bash", displayCmds["default"])
		_, _, _ = runScript("bash", fmt.Sprintf("feh --bg-fill %s", expandHome(cfg.DirWallpaper)))
		return nil
	}

	if _, stderr, err := runScript("bash", cmd); err != nil {
		return fmt.Errorf("set display: %s", stderr)
	}

	_, _, _ = runScript("bash", fmt.Sprintf("feh --bg-fill %s", expandHome(cfg.DirWallpaper)))
	return nil
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
	if _, _, err := runScript("bash", "sudo light -S 48"); err != nil {
		return err
	}
	s.notify("set light to 92%")
	if _, _, err := runScript("bash", "amixer set Master unmute && amixer set Capture cap"); err != nil {
		return err
	}
	if _, _, err := runScript("bash", "amixer set Master 80%"); err != nil {
		return err
	}
	s.notify("set master volume to 80%")
	if _, _, err := runScript("bash", "amixer set Capture 64%"); err != nil {
		return err
	}
	s.notify("set capture volume to 72%")
	if _, _, err := runScript("bash", "xset r rate 158 128"); err != nil {
		return err
	}
	s.notify("set keyboard rate to 158 128")
	return nil
}

func (s *Service) SysKill() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e zsh -c 'ps -ef | fzf --prompt="kill -9 >" --select-1 --exit-0 | awk "{print \$2}" | xargs -r kill -9'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}

func (s *Service) SysOpenTerminal() error {
	term := psl.GetConfig().Svc.DefaultTerminal
	tmpl := `%s -e zsh -lc 'dir=$(fd -t d . "$HOME/share/github" | fzf --prompt="where to open st>" --height 40%% --reverse --preview "ls -lah {}"); if [[ -n "$dir" ]]; then cd "$dir"; fi; exec zsh -i'`
	return startScript("bash", fmt.Sprintf(tmpl, term))
}
