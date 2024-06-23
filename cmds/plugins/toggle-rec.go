package plugins

import (
	"fmt"
	"os"
	"path"
	"time"

	"cmds/sugar"
)

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
				"st  -t %s -c %s -e %s",
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
				"st -g %s -t %s -c %s -e %s",
				sugar.GetGeoForSt(0.74, 0.08, 40, 12),
				WinNameFloatWindow,
				WinNameFloatWindow,
				"ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0",
			),
		)
	}
}
