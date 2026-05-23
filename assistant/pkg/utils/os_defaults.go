package utils

import (
	"fmt"
	"runtime"
)

func GetOSDefault(objType string) (string, error) {
	switch objType {
	case "shell":
		return GetOSDefaultShell()
	case "terminal":
		return GetOSDefaultTerminal()
	case "editor":
		return GetOSDefaultEditor()
	case "browser":
		return GetOSDefaultBrowser()
	default:
		return "", fmt.Errorf("unsupported objType: %q", objType)
	}
}

func GetOSDefaultShell() (string, error) {
	switch runtime.GOOS {
	case OSLinux, OSDarwin:
		return "sh", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func GetOSDefaultTerminal() (string, error) {
	switch runtime.GOOS {
	case OSLinux:
		return "st", nil
	case OSDarwin:
		return "kitty", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func GetOSDefaultEditor() (string, error) {
	switch runtime.GOOS {
	case OSLinux, OSDarwin:
		return "nvim", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func GetOSDefaultBrowser() (string, error) {
	switch runtime.GOOS {
	case OSLinux, OSDarwin:
		return "chrome", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
