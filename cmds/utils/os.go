package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func GetEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

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

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsFile(path string) (isFile bool) {
	return !IsDir(path)
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsDirExists(path string) (exist bool) {
	if Exists(path) && IsDir(path) {
		return true
	}
	return false
}

func IsFileExists(path string) (exist bool) {
	if Exists(path) && !IsDir(path) {
		return true
	}
	return false
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func GetAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func GetPID() int {
	return os.Getpid()
}

func IsRunning(proc string) (isrunning bool) {
	curpid := os.Getpid()
	proc = strings.ReplaceAll(strings.ReplaceAll(proc, "'", ""), `"`, "")
	procs, err := process.Processes()
	if err != nil {
		return false
	}
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if p.Pid == int32(curpid) {
			continue
		}
		if name == proc {
			return true
		}
	}
	for _, p := range procs {
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}
		if p.Pid == int32(curpid) {
			continue
		}
		if strings.Contains(cmdline, proc) {
			return true
		}
	}
	return false
}

func Kill(proc string) {
	curpid := os.Getpid()
	proc = strings.ReplaceAll(strings.ReplaceAll(proc, "'", ""), `"`, "")
	procs, err := process.Processes()
	if err != nil {
		return
	}
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if p.Pid == int32(curpid) {
			continue
		}
		if name == proc {
			p.Kill()
		}
	}
	for _, p := range procs {
		cmdline, err := p.Cmdline()
		if err != nil {
			continue
		}
		if p.Pid == int32(curpid) {
			continue
		}
		if strings.Contains(cmdline, proc) {
			p.Kill()
		}
	}
}
