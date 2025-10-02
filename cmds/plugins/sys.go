package plugins

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"cmds/utils"
)

func SysToggleKeyboardLight() {
	kbdCtlFilePath := "/sys/class/leds/tpacpi::kbd_backlight/brightness"
	brightness, err := utils.GetKeyBoardStatus(kbdCtlFilePath)
	if err != nil {
		utils.Notify(err)
		return
	}
	if brightness == 1 {
		brightness = 0
	} else {
		brightness = 1
	}
	cmd := fmt.Sprintf("sudo sh -c 'echo %d > %s'", brightness, kbdCtlFilePath)
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
}

func SysBlueToothConnect() {
	cmd := "bluetoothctl devices"
	stdout, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		utils.Notify("no bluetooth device was found")
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to'", stdout)
	stdout, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	slice := strings.Split(strings.TrimSpace(stdout), " ")
	if len(slice) != 3 {
		utils.Notify("connect to bluetooth failed")
	}
	deviceid := slice[1]
	cmd = fmt.Sprintf("bluetoothctl connect %s", deviceid)
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	utils.Notify("connect to bluetooth success")
}

func SysBlueToothDisconnect() {
	cmd := "bluetoothctl info | grep 'Device '"
	stdout, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		utils.Notify("no connected bluetooth device found")
		return
	}

	cmd = fmt.Sprintf("echo '%s' | dmenu -p 'disconnect from'", stdout)
	stdout, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}

	slice := strings.Split(strings.TrimSpace(stdout), " ")
	if len(slice) < 2 {
		utils.Notify("disconnect bluetooth failed")
		return
	}
	deviceid := slice[1]

	cmd = fmt.Sprintf("bluetoothctl disconnect %s", deviceid)
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	utils.Notify(fmt.Sprintf("disconnected from %s", deviceid))
}

func SysBlueToothScanAndConnect() {
	devices, err := SysBlueToothScan()
	if err != nil || len(devices) == 0 {
		utils.Notify("没有发现蓝牙设备")
		return
	}

	cmd := exec.Command("dmenu", "-p", "connect bluetooth")
	cmd.Stdin = strings.NewReader(strings.Join(devices, "\n"))
	out, _ := cmd.Output()
	choice := strings.TrimSpace(string(out))
	if choice == "" {
		utils.Notify("未选择设备")
		return
	}

	// parse mac
	parts := strings.Fields(choice)
	if len(parts) < 1 {
		utils.Notify("选择无效")
		return
	}
	mac := parts[0]

	// pair -> trust -> connect
	for _, c := range []string{
		fmt.Sprintf("bluetoothctl pair %s", mac),
		fmt.Sprintf("bluetoothctl trust %s", mac),
		fmt.Sprintf("bluetoothctl connect %s", mac),
	} {
		_, _, err := utils.RunScript("bash", c)
		if err != nil {
			utils.Notify("执行失败: " + c)
			return
		}
	}

	utils.Notify("蓝牙连接成功: " + mac)
}

func SysBlueToothScan() ([]string, error) {
	cmd := exec.Command("bluetoothctl")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// start scan
	_, _ = stdin.Write([]byte("scan on\n"))

	scanner := bufio.NewScanner(stdout)
	found := make(map[string]string)

	re := regexp.MustCompile(`Device\s+([0-9A-F:]{17})\s+(.+)$`)

	// read output for 6 seconds
	timer := time.After(6 * time.Second)
loop:
	for {
		select {
		case <-timer:
			break loop
		default:
			if !scanner.Scan() {
				break loop
			}
			line := scanner.Text()
			// fmt.Println(line)
			if strings.Contains(line, "Device") {
				if m := re.FindStringSubmatch(line); m != nil {
					mac := m[1]
					name := m[2]
					found[mac] = name
					fmt.Printf("MAC=%s, NAME=%s\n", mac, name)
				}
			}
		}
	}

	// stop scan
	_, _ = stdin.Write([]byte("scan off\nexit\n"))
	_ = cmd.Wait()

	var list []string
	for mac, name := range found {
		list = append(list, fmt.Sprintf("%s %s", mac, name))
	}
	return list, nil
}

func SysWifiConnect() {
	cmd := "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq"
	stdout, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to wifi'", stdout)
	stdout, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	essid := strings.TrimSpace(stdout)
	if essid == "" {
		return
	}
	cmd = "dmenu < /dev/null -p 'password'"
	stdout, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	password := strings.TrimSpace(stdout)
	cmd = fmt.Sprintf("nmcli device wifi connect %s password %s", essid, password)
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	utils.Notify("wifi connect success")
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
	content, err := utils.Choose(": ", list)
	if err != nil {
		return
	}
	cmd, ok := SysShortCuts[content]
	if ok {
		_, _, _ = utils.RunScript("bash", cmd)
	}
}

func SysScreen() {
	primaryMonitor := "eDP-1"
	cmd := "xrandr|grep ' connected'|grep -v 'eDP-1'|awk '{print $1}'"
	stdout, _, err := utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		utils.Notify("only have one monitor")
		return
	}
	secondMonitor := stdout
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
	_, _, err = utils.RunScript("bash", cmds["default"])
	if err != nil {
		utils.Notify(err)
		return
	}
	list := make([]string, 0)
	for k := range cmds {
		list = append(list, k)
	}
	choice, err := utils.Choose("screen strategy: ", list)
	if err != nil {
		utils.Notify(err)
		return
	}
	cmd, ok := cmds[choice]
	if !ok {
		utils.Notify("wrong choice")
		return
	}
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
	time.Sleep(10 * time.Millisecond)
	cmd = fmt.Sprintf("feh --bg-fill %s", path.Join(os.Getenv("HOME"), WallPaperPath))
	_, _, err = utils.RunScript("bash", cmd)
	if err != nil {
		utils.Notify(err)
		return
	}
}
