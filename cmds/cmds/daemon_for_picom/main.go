package main

import (
	"cmds/utils"
	"os"
	"time"
)

func main() {
	for {
		if !utils.IsRunning("picom") {
			utils.RunScript("bash", "picom --config "+os.Getenv("HOME")+"/.config/picom/picom.conf")
		}
		time.Sleep(3 * time.Second)
	}
}
