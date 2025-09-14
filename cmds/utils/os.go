package utils

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// -------------------- 环境变量 --------------------
func GetEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// -------------------- 系统信息 --------------------
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

// -------------------- 文件/目录 --------------------
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

// -------------------- 进程 --------------------
func GetPID() int {
	return os.Getpid()
}
