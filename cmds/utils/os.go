package utils

import (
	"errors"
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
	if err == nil {
		return true
	}
	return !errors.Is(err, os.ErrNotExist)
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsDirExists(path string) bool {
	return IsDir(path)
}

func IsFileExists(path string) bool {
	return IsFile(path)
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

func findProcessesByName(proc string) ([]*process.Process, error) {
	curpid := os.Getpid()
	proc = strings.ReplaceAll(strings.ReplaceAll(proc, "'", ""), `"`, "")
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	var result []*process.Process
	for _, p := range procs {
		if p.Pid == int32(curpid) {
			continue
		}
		if name, err := p.Name(); err == nil && name == proc {
			result = append(result, p)
			continue
		}
		if cmdline, err := p.Cmdline(); err == nil && strings.Contains(cmdline, proc) {
			result = append(result, p)
		}
	}
	return result, nil
}

func IsRunning(proc string) bool {
	procs, err := findProcessesByName(proc)
	return err == nil && len(procs) > 0
}

func Kill(proc string) {
	procs, err := findProcessesByName(proc)
	if err != nil {
		return
	}
	for _, p := range procs {
		_ = p.Kill()
	}
}
