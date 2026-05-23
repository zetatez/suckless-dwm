package utils

import (
	"fmt"
	"runtime"
)

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
