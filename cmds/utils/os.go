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

func IsFile(path string) (isFile bool) {
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

func IsDirExists(path string) (exist bool) {
	return IsDir(path)
}

func IsFileExists(path string) (exist bool) {
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

func IsRunning(proc string) (isrunning bool) {
	curpid := os.Getpid()
	proc = strings.ReplaceAll(strings.ReplaceAll(proc, "'", ""), `"`, "")
	procs, err := process.Processes()
	if err != nil {
		return false
	}
	for _, p := range procs {
		if p.Pid == int32(curpid) {
			continue
		}
		if name, err := p.Name(); err == nil && name == proc {
			return true
		}
		if cmdline, err := p.Cmdline(); err == nil && strings.Contains(cmdline, proc) {
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
	killed := map[int32]struct{}{}
	for _, p := range procs {
		if p.Pid == int32(curpid) {
			continue
		}
		if _, ok := killed[p.Pid]; ok {
			continue
		}
		if name, err := p.Name(); err == nil && name == proc {
			_ = p.Kill()
			killed[p.Pid] = struct{}{}
			continue
		}
		if cmdline, err := p.Cmdline(); err == nil && strings.Contains(cmdline, proc) {
			_ = p.Kill()
			killed[p.Pid] = struct{}{}
			continue
		}
	}
}
