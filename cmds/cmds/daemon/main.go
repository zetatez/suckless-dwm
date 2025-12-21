package main

import (
	"cmds/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ProcConfig struct {
	Name    string
	Command string
}

func isRunning(name string) bool {
	selfPID := os.Getpid()
	out, err := exec.Command("pgrep", "-x", name).Output()
	if err != nil {
		return false
	}
	pids := strings.Fields(string(out))
	for _, pidStr := range pids {
		if pidStr == fmt.Sprintf("%d", selfPID) {
			continue
		}
		return true
	}
	return false
}

func main() {
	procs := []ProcConfig{
		{
			Name:    "dwmblocks",
			Command: "dwmblocks",
		},
		{
			Name:    "picom",
			Command: "picom --config " + os.Getenv("HOME") + "/.config/picom/picom.conf &",
		},
		{
			Name:    "dust",
			Command: "dust",
		},
		{
			Name:    "xset",
			Command: "xset r rate 158 128",
		},
		{
			Name:    "clipmenud",
			Command: "clipmenud",
		},
		{
			Name:    "hhkb",
			Command: "hhkb",
		},
	}

	for {
		for _, proc := range procs {
			if !isRunning(proc.Name) {
				go utils.RunScript("bash", proc.Command)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
