package plugins

import (
	"fmt"
	"os"
	"path"
	"time"

	"cmds/sugar"
)

func ToggleRecAudio() {
	cmd := fmt.Sprintf(
		"st -t %s -c %s -e ffmpeg -y -r 60 -f alsa -i default -c:a flac %s",
		WinNameScratchPad,
		WinNameScratchPad,
		path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-audio-%s.flac", time.Now().Local().Format("2006-01-02-15-04-05"))),
	)
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleRecScreen() {
	w, h := sugar.GetScreenSize()
	cmd := fmt.Sprintf(
		"st -t %s -c %s -e ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s",
		WinNameScratchPad,
		WinNameScratchPad,
		w,
		h,
		os.Getenv("DISPLAY"),
		path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-screen-%s.mkv", time.Now().Local().Format("2006-01-02-15-04-05"))),
	)
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleRecWebcam() {
	cmd := fmt.Sprintf(
		"st  -t %s -c %s -e ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s",
		WinNameScratchPad,
		WinNameScratchPad,
		path.Join(os.Getenv("HOME"), fmt.Sprintf("/Videos/rec-webcam-%s.mp4", time.Now().Local().Format("2006-01-02-15-04-05"))),
	)
	switch {
	case sugar.IsRunning("ffmpeg"):
		sugar.Kill("ffmpeg")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}

func ToggleShow() {
	cmd := fmt.Sprintf(
		"st -g %s -t %s -c %s -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
		sugar.GetGeoForSt(0.74, 0.08, 40, 12),
		WinNameFloatWindow,
		WinNameFloatWindow,
	)
	switch {
	case sugar.IsRunning("ffplay"):
		sugar.Kill("ffplay")
	default:
		sugar.NewExecService().RunScriptShell(cmd)
	}
}
