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
	sugar.Toggle(fmt.Sprintf("%s -e abook", sugar.GetOSTerminal()))
}

func ToggleCalendar() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -g %s -t %s -c %s -e nvim +':set laststatus=0' +'Calendar -view=month'", sugar.GetGeoForTerminal(0.84, 0.04, 24, 12), WinNameFloatWindow, WinNameFloatWindow))
	case "kitty":
		sugar.Toggle("kitty -e nvim +':set laststatus=0' +'Calendar -view=month'")
	}
}

func ToggleCalendarSchedulingToday() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -g %s -t %s -c %s -e nvim +':set laststatus=0' +'Calendar -view=day'", sugar.GetGeoForTerminal(0.80, 0.05, 36, 32), WinNameFloatWindow, WinNameFloatWindow))
	case "kitty":
		sugar.Toggle("kitty -e nvim +':set laststatus=0' +'Calendar -view=day'")
	}
}

func ToggleCalendarScheduling() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle("st -t scheduling -c scheduling -e nvim +':set laststatus=0' +'Calendar -view=week'")
	case "kitty":
		sugar.Toggle("kitty -e nvim +':set laststatus=0' +'Calendar -view=week'")
	}
}

func ToggleFlameshot() {
	sugar.Toggle("flameshot gui")
}

func ToggleInkscape() {
	sugar.Toggle("inkscape")
}

func ToggleYazi() {
	sugar.Toggle(fmt.Sprintf("%s -e yazi", sugar.GetOSTerminal()))
}

func ToggleJulia() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e julia", WinNameScratchPad, WinNameScratchPad))
	case "kitty":
		sugar.Toggle(fmt.Sprintf("kitty -T %s -e julia", WinNameScratchPad))
	}
}

func ToggleLazyDocker() {
	sugar.Toggle(fmt.Sprintf("%s -e lazydocker", sugar.GetOSTerminal()))
}

func ToggleLazyGit() {
	sugar.Toggle(fmt.Sprintf("%s -e lazygit", sugar.GetOSTerminal()))
}

func ToggleMusic() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		sugar.Toggle(fmt.Sprintf("%s -e cava", sugar.GetOSTerminal()))
	}()
	time.Sleep(10 * time.Millisecond)
	go func() {
		defer wg.Done()
		sugar.Toggle(fmt.Sprintf("%s -e ncmpcpp", sugar.GetOSTerminal()))
	}()
	wg.Wait()
}

func ToggleMusicNetCloud() {
	sugar.Toggle("netease-cloud-music")
}

func ToggleMutt() {
	sugar.Toggle(fmt.Sprintf("%s -e mutt", sugar.GetOSTerminal()))
}

func ToggleKrita() {
	sugar.Toggle("krita")
}

func TogglePython() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e python -i -c 'import os, sys, datetime, re, json, collections, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt'", WinNameScratchPad, WinNameScratchPad))
	case "kitty":
		sugar.Toggle("kitty -e python -i -c 'import os, sys, datetime, re, json, collections, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt'")
	}
}

func ToggleScala() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e scala", WinNameScratchPad, WinNameScratchPad))
	case "kitty":
		sugar.Toggle("kitty -e scala")
	}
}

func ToggleLua() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -t %s -c %s -e lua", WinNameScratchPad, WinNameScratchPad))
	case "kitty":
		sugar.Toggle("kitty -e lua")
	}
}

func ToggleTTYClock() {
	switch sugar.GetOSTerminal() {
	case "st":
		sugar.Toggle(fmt.Sprintf("st -g %s -t %s -c %s -e tty-clock -s", sugar.GetGeoForTerminal(0.72, 0.04, 53, 8), WinNameFloatWindow, WinNameFloatWindow))
	case "kitty":
		sugar.Toggle("kitty -e tty-clock -s")
	}
}

func ToggleIrssi() {
	sugar.Toggle(fmt.Sprintf("%s -e irssi", sugar.GetOSTerminal()))
}

func ToggleNewsboat() {
	sugar.Toggle(fmt.Sprintf("%s -e newsboat", sugar.GetOSTerminal()))
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

func ToggleTop() {
	sugar.Toggle(fmt.Sprintf("%s -e hop", sugar.GetOSTerminal()))
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
	stdout, _, err := sugar.NewExecService().RunScript("bash", cmd)
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
	_, _, err = sugar.NewExecService().RunScript("bash", cmd)
	if err != nil {
		sugar.Notify(err)
		return
	}
}

func ToggleXournal() {
	sugar.Toggle("xournalpp")
}

func ToggleRecAudio() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		switch sugar.GetOSTerminal() {
		case "st":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("st -t %s -c %s -e %s", WinNameScratchPad, WinNameScratchPad, fmt.Sprintf("ffmpeg -y -r 60 -f alsa -i default -c:a flac %s", path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-audio-%s.flac", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		case "kitty":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("kitty -t %s -e %s", WinNameScratchPad, fmt.Sprintf("ffmpeg -y -r 60 -f alsa -i default -c:a flac %s", path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-audio-%s.flac", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		}

	}
}

func ToggleRecScreen() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		w, h := sugar.GetScreenSize()
		switch sugar.GetOSTerminal() {
		case "st":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("st -t %s -c %s -e %s ", WinNameScratchPad, WinNameScratchPad, fmt.Sprintf("ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s", w, h, os.Getenv("DISPLAY"), path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-screen-%s.mkv", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		case "kitty":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("kitty -T %s -e %s ", WinNameScratchPad, fmt.Sprintf("ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s", w, h, os.Getenv("DISPLAY"), path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-screen-%s.mkv", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		}
	}
}

func ToggleRecWebcam() {
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		switch sugar.GetOSTerminal() {
		case "st":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("st -t %s -c %s -e %s", WinNameScratchPad, WinNameScratchPad, fmt.Sprintf("ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s", path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-webcam-%s.mp4", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		case "kitty":
			sugar.NewExecService().RunScript("bash", fmt.Sprintf("kitty -T %s -e %s", WinNameScratchPad, fmt.Sprintf("ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s", path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-webcam-%s.mp4", time.Now().Local().Format("2006-01-02-15-04-05"))))))
		}
	}
}

func ToggleShow() {
	switch {
	case sugar.IsRunning("ffplay"):
		sugar.Kill("ffplay")
	default:
		sugar.NewExecService().RunScript("bash",
			fmt.Sprintf("%s -e %s", sugar.GetOSTerminal(), "ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0"),
		)
	}
}
