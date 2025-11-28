package plugins

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"cmds/utils"
)

func ToggleCalendar() {
	utils.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e nvim +':Calendar -view=month'",
			utils.GetGeoForTerminal(0.84, 0.04, 24, 12),
			winClassFloat,
			winClassFloat,
		),
	)
}

func ToggleCalendarSchedulingToday() {
	utils.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e nvim +':Calendar -view=day'",
			utils.GetGeoForTerminal(0.80, 0.05, 36, 32),
			winClassFloat,
			winClassFloat,
		),
	)
}

func ToggleCalendarScheduling() {
	utils.Toggle("st -t scheduling -c scheduling -e nvim +':Calendar -view=week'")
}

func ToggleFlameshot() {
	utils.Toggle("flameshot gui")
}

func ToggleYazi() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e yazi",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleJulia() {
	utils.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e julia",
			WinClassScratchPad,
			WinClassScratchPad,
		),
	)
}

func ToggleLazyDocker() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e lazydocker",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleLazyGit() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e lazygit",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleMusic() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		utils.Toggle(
			fmt.Sprintf(
				"%s -e cava",
				utils.GetOSDefaultTerminal(),
			),
		)
	}()
	time.Sleep(10 * time.Millisecond)
	go func() {
		defer wg.Done()
		utils.Toggle(
			fmt.Sprintf(
				"%s -e ncmpcpp",
				utils.GetOSDefaultTerminal(),
			),
		)
	}()
	wg.Wait()
}

func ToggleNeteaseCloudMusic() {
	utils.Toggle("netease-cloud-music")
}

func ToggleMutt() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e mutt",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func TogglePython() {
	utils.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e python -i -c 'import os, sys, datetime, re, json, collections, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt; print(dir())'",
			WinClassScratchPad,
			WinClassScratchPad,
		),
	)
}

func ToggleScala() {
	utils.Toggle(
		fmt.Sprintf(
			"st -t %s -c %s -e scala",
			WinClassScratchPad,
			WinClassScratchPad,
		),
	)
}

func ToggleTTYClock() {
	utils.Toggle(
		fmt.Sprintf(
			"st -g %s -t %s -c %s -e tty-clock -s",
			utils.GetGeoForTerminal(0.72, 0.04, 53, 8),
			winClassFloat,
			winClassFloat,
		),
	)
}

func ToggleIrssi() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e irssi",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleNewsboat() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e newsboat",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleScreenKey() {
	utils.Toggle("screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow")
}

func ToggleTop() {
	utils.Toggle(
		fmt.Sprintf(
			"%s -e btop",
			utils.GetOSDefaultTerminal(),
		),
	)
}

func ToggleClipmenu() {
	utils.Toggle("sh -c clipmenu")
}

func TogglePassmenu() {
	utils.Toggle("passmenu")
}

func ToggleRecAudio() {
	switch {
	case utils.IsRunning("ffmpeg"):
		utils.Kill("ffmpeg")
	default:
		_, _, _ = utils.RunScript(
			"bash",
			fmt.Sprintf(
				"st -t %s -c %s -e ffmpeg -y -r 60 -f alsa -i default -c:a flac %s",
				WinClassScratchPad,
				WinClassScratchPad,
				path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-audio-%s.flac", time.Now().Local().Format("2006-01-02-15-04-05"))),
			),
		)
	}
}

func ToggleRecScreen() {
	switch {
	case utils.IsRunning("ffmpeg"):
		utils.Kill("ffmpeg")
	default:
		w, h := utils.GetScreenSize()
		_, _, _ = utils.RunScript(
			"bash",
			fmt.Sprintf(
				"st -t %s -c %s -e ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s",
				WinClassScratchPad,
				WinClassScratchPad,
				w,
				h,
				os.Getenv("DISPLAY"),
				path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-screen-%s.mkv", time.Now().Local().Format("2006-01-02-15-04-05"))),
			),
		)
	}
}

func ToggleRecWebcam() {
	switch {
	case utils.IsRunning("ffmpeg"):
		utils.Kill("ffmpeg")
	default:
		_, _, _ = utils.RunScript(
			"bash",
			fmt.Sprintf(
				"st -t %s -c %s -e ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s",
				WinClassScratchPad,
				WinClassScratchPad,
				path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-webcam-%s.mp4", time.Now().Local().Format("2006-01-02-15-04-05"))),
			),
		)
	}
}

func ToggleShow() {
	switch {
	case utils.IsRunning("ffplay"):
		utils.Kill("ffplay")
	default:
		_, _, _ = utils.RunScript(
			"bash",
			"st -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
		)
	}
}
