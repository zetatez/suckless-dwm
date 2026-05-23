package utils

import (
	"fmt"
	"runtime"
)

func Shutdown() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("systemctl", "poweroff")
	case OSDarwin:
		return runCmd("osascript", "-e", `tell application "System Events" to shut down`)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Reboot() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("systemctl", "reboot")
	case OSDarwin:
		return runCmd("osascript", "-e", `tell application "System Events" to restart`)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Suspend() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("systemctl", "suspend")
	case OSDarwin:
		return runCmd("pmset", "sleepnow")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Lock() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("slock")
	case OSDarwin:
		return runCmd(
			"osascript",
			"-e",
			`tell application "System Events" to keystroke "q" using {control down, command down}`,
		)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
