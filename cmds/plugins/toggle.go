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
	_, _, err = sugar.NewExecService().RunScriptShell(cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
	sugar.Notify("connect to bluetooth success")
}

func ToggleCalendar() {
	sugar.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e nvim +':set laststatus=0' +'Calendar -view=month'",
			sugar.GetGeoForSt(0.84, 0.04, 24, 12),
			WinNameFloatWindow,
			WinNameFloatWindow,
		),
	)
}

func ToggleCalendarSchedulingToday() {
	sugar.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e nvim +':set laststatus=0' +'Calendar -view=day'",
			sugar.GetGeoForSt(0.80, 0.05, 36, 32),
			WinNameFloatWindow,
			WinNameFloatWindow,
		),
	)
}

func ToggleCalendarScheduling() {
	sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=week'")
}

func ToggleDarkTable() {
	sugar.Toggle("darktable")
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
	sugar.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e julia",
			WinNameScratchPad,
			WinNameScratchPad,
		),
	)
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

func ToggleLazyGit() {
	sugar.Toggle("st -e lazygit")
}

func ToggleMusic() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		sugar.Toggle("st -e cava")
	}()
	time.Sleep(10 * time.Millisecond)
	go func() {
		defer wg.Done()
		sugar.Toggle("st -e ncmpcpp")
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
	sugar.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e python -i -c 'import os, sys, datetime, re, json, collections, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt'",
			WinNameScratchPad,
			WinNameScratchPad,
		),
	)
}

func ToggleScala() {
	sugar.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e scala",
			WinNameScratchPad,
			WinNameScratchPad,
		),
	)
}

func ToggleLua() {
	sugar.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e lua",
			WinNameScratchPad,
			WinNameScratchPad,
		),
	)
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

func ToggleObsidian() {
	sugar.Toggle("obsidian")
}

func ToggleSysShortcuts() {
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

func ToggleTermius() {
	sugar.Toggle("termius")
}

func ToggleRecAudio() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"st -t %s -c %s -e %s",
				WinNameScratchPad,
				WinNameScratchPad,
				fmt.Sprintf(
					"ffmpeg -y -r 60 -f alsa -i default -c:a flac %s",
					path.Join(
						os.Getenv("HOME"),
						fmt.Sprintf(
							"/Videos/rec-audio-%s.flac",
							time.Now().Local().Format("2006-01-02-15-04-05"),
						),
					),
				),
			),
		)
	}
}

func ToggleRecScreen() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		w, h := sugar.GetScreenSize()
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"st -t %s -c %s -e %s ",
				WinNameScratchPad,
				WinNameScratchPad,
				fmt.Sprintf(
					"ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s",
					w,
					h,
					os.Getenv("DISPLAY"),
					path.Join(
						os.Getenv("HOME"),
						fmt.Sprintf(
							"/Videos/rec-screen-%s.mkv",
							time.Now().Local().Format("2006-01-02-15-04-05"),
						),
					),
				),
			),
		)
	}
}

func ToggleRecWebcam() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"st -t %s -c %s -e %s",
				WinNameScratchPad,
				WinNameScratchPad,
				fmt.Sprintf(
					"ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s",
					path.Join(
						os.Getenv("HOME"),
						fmt.Sprintf(
							"/Videos/rec-webcam-%s.mp4",
							time.Now().Local().Format("2006-01-02-15-04-05"),
						),
					),
				),
			),
		)
	}
}

func ToggleShow() {
	switch {
	case sugar.IsRunning("ffplay"):
		sugar.Kill("ffplay")
	default:
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"st -e %s",
				"ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
			),
		)
	}
}
