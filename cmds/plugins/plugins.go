package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"cmds/sugar"

	"golang.design/x/clipboard"
)

const (
	WallPaperPath    = "Pictures/wallpapers"
	DefaultWallpaper = "0101.jpg"
	ScratchPad       = "scratchpad"
	FloatWindow      = "00001011"
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

func ToggleCalendarDay() {
	sugar.Toggle(
		fmt.Sprintf("st -g %s -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=day'",
			sugar.GetGeoForSt(0.80, 0.05, 36, 32),
		),
	)
}

func ToggleCalendarWeek() {
	sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=week'")
}

func ToggleChrome() {
	sugar.Toggle("chrome --proxy-server=socks5://127.0.0.1:7891")
}

func ToggleDiary() {
	// TODO: <13:26:58 2024-01-14: Dionysus>:
}

func ToggleFlameshot() {
	sugar.Toggle("flameshot gui")
}

func ToggleInkscape() {
	LaunchApp("inkscape")
}

func ToggleJoshuto() {
	sugar.Toggle("st -e joshuto")
}

func ToggleJulia() {
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e julia", ScratchPad, ScratchPad))
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
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e python", ScratchPad, ScratchPad))
}

func ToggleRecAudio() {
	cmd := fmt.Sprintf(
		"st  -t %s -c %s -e ffmpeg -y -r 60 -f alsa -i default -c:a flac %s",
		ScratchPad,
		ScratchPad,
		path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-a-%s.flac", time.Now().Local().Format("2006-01-02-15-04-05"))),
	)
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleRecVideo() {
	w, h := sugar.GetScreenSize()
	cmd := fmt.Sprintf(
		"st  -t %s -c %s -e ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s",
		ScratchPad,
		ScratchPad,
		w,
		h,
		os.Getenv("DISPLAY"),
		path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-v-a-%s.mkv", time.Now().Local().Format("2006-01-02-15-04-05"))),
	)
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleIrssi() {
	sugar.Toggle("st -e irssi")
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
	choice, err := sugar.Choose(list)
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
	cmd = fmt.Sprintf("feh --bg-file %s", path.Join(os.Getenv("HOME"), WallPaperPath, DefaultWallpaper))
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func ToggleScreenKey() {
	sugar.Toggle("screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow")
}

func ToggleShow() {
	LaunchApp(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
			sugar.GetGeoForSt(0.74, 0.08, 40, 12),
			FloatWindow,
			FloatWindow,
		),
	)
}

func ToggleSublime() {
	sugar.Toggle("subl")
}

func ToggleSysShortcuts() {
	SysShortCuts := map[string]string{
		"󰒲 suspend":     "systemctl suspend",
		" poweroff":    "systemctl poweroff",
		"ﰇ reboot":      "systemctl reboot",
		"󰷛 slock":       "slock & sleep 0.5 & xset dpms force off",
		"󰶐 off-display": "sleep .5; xset dpms force off",
	}
	list := []string{}
	for k := range SysShortCuts {
		list = append(list, k)
	}
	content, err := sugar.Choose(list)
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
	LaunchApp("clipmenu")
}

func TogglePassmenu() {
	LaunchApp("passmenu")
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

func UmountXYZ() {
	cmd := "echo '/x\n/y\n/z'|dmenu -p 'umount'"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	choice := stdout
	cmd = fmt.Sprintf("sudo umount %s", choice)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("umount success")
}

func WifiConnect() {
	cmd := "nmcli device wifi list|sed '1d'|sed '/--/ d'|awk '{print $2}'|sort|uniq"
	stdout, _, err := sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	cmd = fmt.Sprintf("echo '%s'|dmenu -p 'connect to wifi'", stdout)
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	essid := strings.TrimSpace(stdout)
	if essid == "" {
		return
	}
	cmd = "dmenu < /dev/null -p 'password'"
	stdout, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	password := strings.TrimSpace(stdout)
	cmd = fmt.Sprintf("nmcli device wifi connect %s password %s", essid, password)
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("wifi connect success")
}

func CurrentDatetime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := time.Now().Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func CurrentUnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := fmt.Sprintf("%d", time.Now().Unix())
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func HandleCopied() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	content := strings.TrimSpace(string(text))
	switch {
	case sugar.Exists(content) && sugar.IsFile(content):
		sugar.Lazy("open", content)
		return
	case sugar.IsUrl(content):
		Website(content)()
		return
	default:
		SearchFromWeb(content)
	}
}

func LaunchApp(cmd string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func Website(url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 %s", url),
		)
	}
}

func SearchFromWeb(content string) {
	sugar.NewExecService().RunScriptShell(
		fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 https://cn.bing.com/search?q=%s", content),
	)
}

func TransferDatetime2UnixSec() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	t, err := time.Parse(time.DateTime, strings.TrimSpace(string(text)))
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := fmt.Sprintf("%d", t.Unix())
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func TransferUnixSec2Datetime() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	unix, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := time.Unix(unix, 0).Format(time.DateTime)
	sugar.Notify(fmt.Sprintf("tranfer success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func FormatJson() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)

	doc := map[string]interface{}{}
	err = json.Unmarshal(text, &doc)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, formatedText)
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}

func FormatSql() {
	err := clipboard.Init()
	if err != nil {
		sugar.Notify(err)
		return
	}
	text := clipboard.Read(clipboard.FmtText)
	cmd := `
import sqlparse
print(sqlparse.format("""%s""", reindent=True, indent=2, keyword_case='upper'))
`
	cmd = fmt.Sprintf(cmd, string(text))
	stdout, _, err := sugar.NewExecService().RunScriptPython(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	formatedText := stdout
	sugar.Notify(fmt.Sprintf("format success: \n%s", formatedText))
	changed := clipboard.Write(clipboard.FmtText, []byte(formatedText))
	select {
	case <-changed:
		sugar.Notify("previous clipboard expired")
	}
}
