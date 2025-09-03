package main

import (
	"cmds/sugar"
	"os"
	"time"
)

func main() {
	for {
		if !sugar.IsRunning("picom") {
			sugar.NewExecService().RunScript("bash", "picom --config "+os.Getenv("HOME")+"/.config/picom/picom.conf")
		}
		time.Sleep(1 * time.Second)
	}
}
