package plugins

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"cmds/sugar"
)

func ToggleAddressbook() {
	sugar.Toggle("st -e abook")
}

func ToggleBlueTooth() {
	cmd := "bluetoothctl devices"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		sugar.Notify("no bluetooth device was found")
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to'", stdout)
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	slice := strings.Split(strings.TrimSpace(stdout), " ")
	if len(slice) != 3 {
		sugar.Notify("connect to bluetooth failed")
	}
	deviceid := slice[1]
	cmd = fmt.Sprintf("bluetoothctl connect %s", deviceid)
	fmt.Println(cmd)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("connect to bluetooth success")
}

func ToggleCalendarTodaySchedule() {
	sugar.Toggle(
		fmt.Sprintf("st -g %s -t %s -c %s -e nvim +':set laststatus=0' +'Calendar -view=day'",
			sugar.GetGeoForSt(0.80, 0.05, 36, 32),
			WinNameFloat,
			WinNameFloat,
		),
	)
}

func ToggleCalendarScheduling() {
	sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=week'")
}

func ToggleChrome() {
	sugar.Toggle("chrome --proxy-server=socks5://127.0.0.1:7891")
}

func ToggleEdge() {
	sugar.Toggle("edge --proxy-server=socks5://127.0.0.1:7891")
}

func ToggleFlameshot() {
	sugar.Toggle("flameshot gui")
}

func ToggleInkscape() {
	sugar.Toggle("inkscape")
}

func ToggleJoshuto() {
	sugar.Toggle("st -e joshuto")
}

func ToggleJulia() {
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e julia", WinNameScratchPad, WinNameScratchPad))
}

func ToggleKeyboardLight() {
	kbPath := "/sys/class/leds/tpacpi::kbd_backlight/brightness"
	brightness, err := sugar.GetKeyBoardStatus(kbPath)
	if err != nil {
		sugar.Notify(err)
		return
	}
	if brightness == 1 {
		brightness = 0
	} else {
		brightness = 1
	}
	cmd := fmt.Sprintf("sudo sh -c 'echo %d > %s'", brightness, kbPath)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func ToggleLazyDocker() {
	sugar.Toggle("st -e lazydocker")
}

func ToggleMusic() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		sugar.Toggle(fmt.Sprintf("st -g %s -t cava -c cava -e cava", sugar.GetGeoForSt(0.74, 0.08, 40, 12)))
	}()
	time.Sleep(10 * time.Millisecond)
	go func() {
		defer wg.Done()
		sugar.Toggle(fmt.Sprintf("st -g %s -t music -c music -e ncmpcpp", sugar.GetGeoForSt(0.52, 0.08, 40, 12)))
	}()
	wg.Wait()
}

func ToggleMusicNetCloud() {
	sugar.Toggle("netease-cloud-music")
}

func ToggleMutt() {
	sugar.Toggle("st -e mutt")
}

func ToggleKrita() {
	sugar.Toggle("krita")
}

func TogglePython() {
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e python", WinNameScratchPad, WinNameScratchPad))
}

func ToggleScala() {
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e scala", WinNameScratchPad, WinNameScratchPad))
}

func ToggleLua() {
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e lua", WinNameScratchPad, WinNameScratchPad))
}

func ToggleIrssi() {
	sugar.Toggle("st -e irssi")
}

func ToggleNewsboat() {
	sugar.Toggle("st -e newsboat")
}

func ToggleScreen() {
	primaryMonitor := "eDP-1"
	secondMonitor := "eDP-1"
	cmd := "xrandr|grep ' connected'|grep -v 'eDP-1'|awk '{print $1}'"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		sugar.Notify("have no second monitor")
		return
	}
	secondMonitor = stdout
	cmds := map[string]string{
		"defualt":              fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"clone":                fmt.Sprintf("xrandr --output %s --mode 1920x1080", secondMonitor),
		"monitor only":         fmt.Sprintf("xrandr --output %s --auto --output %s --off", secondMonitor, primaryMonitor),
		"laptop only":          fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"left of":              fmt.Sprintf("xrandr --output %s --auto --left-of %s --auto", secondMonitor, primaryMonitor),
		"right of":             fmt.Sprintf("xrandr --output %s --auto --right-of %s --auto", secondMonitor, primaryMonitor),
		"above":                fmt.Sprintf("xrandr --output %s --auto --above %s --auto", secondMonitor, primaryMonitor),
		"below":                fmt.Sprintf("xrandr --output %s --auto --below %s --auto", secondMonitor, primaryMonitor),
		"roate left left-of":   fmt.Sprintf("xrandr --output %s --auto --rotate left --left-of %s --auto", secondMonitor, primaryMonitor),
		"roate right left-of":  fmt.Sprintf("xrandr --output %s --auto --rotate right --left-of %s --auto", secondMonitor, primaryMonitor),
		"roate left right-of":  fmt.Sprintf("xrandr --output %s --auto --rotate left --right-of %s --auto", secondMonitor, primaryMonitor),
		"roate right right-of": fmt.Sprintf("xrandr --output %s --auto --rotate right --right-of %s --auto", secondMonitor, primaryMonitor),
	}
	_, _, err = sugar.NewExecService().RunScriptShell(cmds["default"])
	if err != nil {
		sugar.Notify(err)
		return
	}
	list := make([]string, 0)
	for k := range cmds {
		list = append(list, k)
	}
	choice, err := sugar.Choose("screen strategy: ", list)
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd, ok := cmds[choice]
	if !ok {
		sugar.Notify("wrong choice")
		return
	}
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	time.Sleep(10 * time.Millisecond)
	cmd = fmt.Sprintf("feh --bg-fill %s", path.Join(os.Getenv("HOME"), WallPaperPath, DefaultWallpaper))
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func ToggleScreenKey() {
	sugar.Toggle("screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow")
}

func ToggleSublime() {
	sugar.Toggle("subl")
}

func ToggleSysShortcuts() {
	SysShortCuts := map[string]string{
		"󰒲  suspend":     "systemctl suspend",
		"  poweroff":    "systemctl poweroff",
		"ﰇ  reboot":      "systemctl reboot",
		"󰷛  slock":       "slock & sleep 0.5 & xset dpms force off",
		"󰶐  off-display": "sleep .5; xset dpms force off",
	}
	list := []string{}
	for k := range SysShortCuts {
		list = append(list, k)
	}
	content, err := sugar.Choose(": ", list)
	if err != nil {
		return
	}
	cmd, ok := SysShortCuts[content]
	if ok {
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleTop() {
	sugar.Toggle("st -e top")
}

func ToggleWallpaper() {
	sugar.Toggle(
		fmt.Sprintf("feh --bg-fill --recursive --randomize %s", path.Join(os.Getenv("HOME"), WallPaperPath)),
	)
}

func ToggleWechat() {
	sugar.Toggle("st -e wechat-uos")
}

func ToggleClipmenu() {
	sugar.Toggle("sh -c clipmenu")
}

func TogglePassmenu() {
	sugar.Toggle("passmenu")
}

func ToggleRedShift() {
	cmd := "systemctl --user status redshift.service"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	switch {
	case strings.Contains(stdout, "running"):
		cmd = "systemctl --user stop redshift.service"
	default:
		cmd = "systemctl --user start redshift.service"
	}
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func ToggleXournal() {
	sugar.Toggle("xournalpp")
}
