package plugins

import (
	"fmt"
	"os"
	"path"
	"time"

	"cmds/sugar"
)

const (
	WallPaperPath = "/home/dionysus/Pictures/wallpapers"
	ScratchPad    = "scratchpad"
	FloatWindow   = "00001011"
)

func ToggleAddressbook() {
	sugar.Toggle("st -e abook")
}

func ToggleBlueTooth() {
	// sugar.Toggle("st -e abook")
}

func ToggleCalendarDay() {
	sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=day'")
}

func ToggleCalendarWeek() {
	sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=week'")
}

func ToggleChrome() {
	sugar.Toggle("chrome")
}

func ToggleChromeWithProxy() {
	sugar.Toggle("chrome --proxy-server=socks5://127.0.0.1:7891")
}

func ToggleDiary() {
	// TODO: <13:26:58 2024-01-14: Dionysus>:
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
	sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e julia", ScratchPad, ScratchPad))
}

func ToggleKeyboardLight() {
	// TODO: <13:30:37 2024-01-14: Dionysus>:
}

func ToggleLazyDocker() {
	sugar.Toggle("st -e lazydocker")
}

func ToggleMusic() {
	// TODO: <13:31:58 2024-01-14: Dionysus>:
	sugar.Toggle("st -g 40x12+100+200 -t cava -c cava -e cava &")
	sugar.Toggle("st -g 40x12+100+200 -t music -c music -e ncmpcpp &")
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
	cmd := fmt.Sprintf(
		"st  -t %s -c %s -e ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s",
		ScratchPad,
		ScratchPad,
		2056, // todo: replace
		1600, // todo: replace
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

func ToggleRedShift() {
	// TODO: <13:36:01 2024-01-14: Dionysus>:
}

func ToggleIrssi() {
	sugar.Toggle("st -e irssi")
}

func ToggleScreen() {
	// TODO: <13:37:18 2024-01-14: Dionysus>:
}

func ToggleScreenKey() {
	sugar.Toggle("screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow")
}

func ToggleShow() {
	sugar.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
			"40x12+600+600",
			FloatWindow,
			FloatWindow,
		),
	)
}

func ToggleSublime() {
	switch {
	case sugar.IsRunning("subl"):
		sugar.Kill("subl")
	default:
		sugar.NewExecService().RunScriptShell("sublime_text")
	}
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
	sugar.Toggle("st -e btop")
}

func ToggleWallpaper() {
	sugar.Toggle(fmt.Sprintf("feh --bg-fill --recursive --randomize %s", WallPaperPath))
}

func ToggleWechat() {
	sugar.Toggle("st -e wechat-uos")
}
