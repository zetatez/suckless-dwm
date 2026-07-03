package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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

// Volume sets output volume (0-100).
func Volume(level int) error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("amixer", "set", "Master", fmt.Sprintf("%d%%", level))
	case OSDarwin:
		return runCmd("osascript", "-e", fmt.Sprintf("set volume output volume %d", level))
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// Mute mutes the output.
func Mute() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("amixer", "set", "Master", "mute")
	case OSDarwin:
		return runCmd("osascript", "-e", "set volume output muted true")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// Unmute unmutes the output.
func Unmute() error {
	switch runtime.GOOS {
	case OSLinux:
		return runCmd("amixer", "set", "Master", "unmute")
	case OSDarwin:
		return runCmd("osascript", "-e", "set volume output muted false")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
