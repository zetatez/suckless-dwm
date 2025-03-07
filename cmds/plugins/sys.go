package plugins

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cmds/sugar"
)

func SysToggleKeyboardLight() {
	kbdCtlFilePath := "/sys/class/leds/tpacpi::kbd_backlight/brightness"
	brightness, err := sugar.GetKeyBoardStatus(kbdCtlFilePath)
	if err != nil {
		sugar.Notify(err)
		return
	}
	if brightness == 1 {
		brightness = 0
	} else {
		brightness = 1
	}
	cmd := fmt.Sprintf("sudo sh -c 'echo %d > %s'", brightness, kbdCtlFilePath)
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func SysBlueTooth() {
	cmd := "bluetoothctl devices"
	stdout, _, err := sugar.NewExecService().RunScript("bash", cmd)
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
	stdout, _, err = sugar.NewExecService().RunScript("bash", cmd)
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
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("connect to bluetooth success")
}

func SysWifiConnect() {
	cmd := "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq"
	stdout, _, err := sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to wifi'", stdout)
	stdout, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	essid := strings.TrimSpace(stdout)
	if essid == "" {
		return
	}
	cmd = "dmenu < /dev/null -p 'password'"
	stdout, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	password := strings.TrimSpace(stdout)
	cmd = fmt.Sprintf("nmcli device wifi connect %s password %s", essid, password)
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("wifi connect success")
}

func SysShortcuts() {
	SysShortCuts := map[string]string{
		"󰒲  suspend":     "systemctl suspend",
		"  poweroff":    "systemctl poweroff",
		"ﰇ  reboot":      "systemctl reboot",
		"󰶐  off-display": "sleep .5; xset dpms force off",
		"󰷛  slock":       "slock",
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
		sugar.NewExecService().RunScript("bash", cmd)
	}
}

func SysScreen() {
	primaryMonitor, secondMonitor := "eDP1", "eDP1"
	cmd := "xrandr|grep ' connected'|grep -v 'eDP1'|awk '{print $1}'"
	stdout, _, err := sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		sugar.Notify("only have one monitor")
		return
	}
	secondMonitor = stdout
	cmds := map[string]string{
		"defualt":                fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"clone":                  fmt.Sprintf("xrandr --output %s --mode 1920x1080", secondMonitor),
		"primary only":           fmt.Sprintf("xrandr --output %s --auto --output %s --off", primaryMonitor, secondMonitor),
		"second  only":           fmt.Sprintf("xrandr --output %s --auto --output %s --off", secondMonitor, primaryMonitor),
		"left  of":               fmt.Sprintf("xrandr --output %s --auto --left-of %s --auto", secondMonitor, primaryMonitor),
		"right of":               fmt.Sprintf("xrandr --output %s --auto --right-of %s --auto", secondMonitor, primaryMonitor),
		"above":                  fmt.Sprintf("xrandr --output %s --auto --above %s --auto", secondMonitor, primaryMonitor),
		"below":                  fmt.Sprintf("xrandr --output %s --auto --below %s --auto", secondMonitor, primaryMonitor),
		"roate left  & left-of":  fmt.Sprintf("xrandr --output %s --auto --rotate left  --left-of %s --auto", secondMonitor, primaryMonitor),
		"roate right & left-of":  fmt.Sprintf("xrandr --output %s --auto --rotate right --left-of %s --auto", secondMonitor, primaryMonitor),
		"roate left  & right-of": fmt.Sprintf("xrandr --output %s --auto --rotate left  --right-of %s --auto", secondMonitor, primaryMonitor),
		"roate right & right-of": fmt.Sprintf("xrandr --output %s --auto --rotate right --right-of %s --auto", secondMonitor, primaryMonitor),
	}
	_, _, err = sugar.NewExecService().RunScript("bash", cmds["default"])
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
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	time.Sleep(10 * time.Millisecond)
	cmd = fmt.Sprintf("feh --bg-fill %s", path.Join(os.Getenv("HOME"), WallPaperPath))
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}
