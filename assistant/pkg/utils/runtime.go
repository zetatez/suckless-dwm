package utils

import (
	"os"
	"os/user"
	"runtime"
)

func GetOSType() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

func GetHostname() (string, error) {
	return os.Hostname()
}

func GetCurrentUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func GetPID() int {
	return os.Getpid()
}
