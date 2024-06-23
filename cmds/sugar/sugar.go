package sugar

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/shirou/gopsutil/process"
)

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

func Exists(path string) (exist bool) {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) (isDir bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func CopyDir(src string, dst string) (err error) {
	if strings.TrimSpace(src) == strings.TrimSpace(dst) {
		return fmt.Errorf("src path %s is equal to dst path %s", src, dst)
	}

	if !IsDirExists(src) {
		return fmt.Errorf("src path %s is not exist", src)
	}

	if !IsDirExists(dst) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}
	}

	absSrc, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	err = filepath.Walk(
		absSrc,
		func(s string, info os.FileInfo, err error) error {
			if s == absSrc {
				return nil
			}
			if info == nil {
				return err
			}
			d := strings.ReplaceAll(s, absSrc, absDst)
			if info.IsDir() {
				if !IsDirExists(d) {
					if err = os.MkdirAll(d, os.ModePerm); err != nil {
						return err
					}
				}
			} else {
				err = CopyFile(s, d)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)

	return err
}

func CopyFile(src string, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	sinfo, err := s.Stat()
	if err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, sinfo.Mode())
	if err != nil {
		return err
	}
	return nil
}

func IsFile(path string) (isFile bool) {
	return !IsDir(path)
}

func Choose(prompt string, list []string) (item string, err error) {
	script := fmt.Sprintf(
		// "echo '%s'|rofi -show -dmenu -p '%s'",
		"echo '%s'|dmenu -p '%s'",
		strings.Join(list, "\n"),
		prompt,
	)
	stdout, _, err := NewExecService().RunScriptShell(script)
	if err != nil {
		return "", err
	}
	item = strings.TrimSpace(stdout)
	return item, nil
}

func GetInput(prompt string) (input string, err error) {
	script := fmt.Sprintf(
		// "dmenu < /dev/null -p '%s'",           // cause: dmenu can not input ch
		"rofi -show -dmenu < /dev/null -p '%s'",
		prompt,
	)
	stdout, _, err := NewExecService().RunScriptShell(script)
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(stdout)
	return input, nil
}

func Lazy(option string, filepath string) {
	switch option {
	case "view", "open", "exec", "copy", "rename", "delete":
		NewExecService().RunScriptShell(
			fmt.Sprintf("st -e lazy -o %s -f %s &", option, filepath),
		)
	default:
		return
	}
}

func IsUrl(content string) (isUrl bool) {
	r := regexp.MustCompile("^(http|www|file).*")
	return r.Match([]byte(content))
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

func Toggle(proc string) {
	if IsRunning(proc) {
		Kill(proc)
	} else {
		NewExecService().RunScriptShell(proc)
	}
}

func Notify(msg ...interface{}) {
	NewExecService().RunScriptShell(fmt.Sprintf("notify-send '%v'", msg))
}

func GetKeyBoardStatus(kbPath string) (brightness int64, err error) {
	brightnessStr, err := os.ReadFile(kbPath)
	if err != nil {
		return 0, err
	}
	brightness, err = strconv.ParseInt(strings.TrimSpace(string(brightnessStr)), 10, 64)
	if err != nil {
		return 0, err
	}
	return brightness, nil
}

func GetScreenSize() (width int, height int) {
	return robotgo.GetScreenSize()
}

func GetPosition(xr float64, yr float64) (x, y int) {
	width, height := GetScreenSize()
	return int(float64(width) * xr), int(float64(height) * yr)
}

func GetGeoForSt(xr float64, yr float64, w int, h int) (geo string) {
	x, y := GetPosition(xr, yr)
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}

func GetGeoCenterForSt(wr float64, hr float64) (geo string) {
	width, height := GetScreenSize()
	Notify(width, height)
	w := int(float64(width) * wr)
	h := int(float64(height) * hr)
	x := (width - w) / 2
	y := (height - h) / 2
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}

func GetKnownHosts() (knownHosts []string, err error) {
	knownHosts = []string{}
	set := map[string]bool{}
	b, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".ssh/known_hosts"))
	if err != nil {
		return knownHosts, err
	}
	str := string(b)
	slice1 := strings.Split(str, "\n")
	for _, x := range slice1 {
		slice2 := strings.Split(x, " ")
		if len(slice2) != 3 {
			continue
		}
		host := strings.TrimSpace(slice2[0])
		if len(host) == 0 {
			continue
		}
		set[host] = true
	}
	for k := range set {
		knownHosts = append(knownHosts, k)
	}
	sort.Strings(knownHosts)
	return knownHosts, nil
}

func SSH(host string, port int, user string, password string) (err error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	cmd := fmt.Sprintf(
		"st -e %s -c '%s'",
		shell,
		fmt.Sprintf(
			`sshpass -p "%s" ssh -o "StrictHostKeyChecking no" -p %d %s@%s`,
			password,
			port,
			user,
			host,
		),
	)
	_, _, err = NewExecService().RunScriptShell(cmd)
	if err != nil {
		return err
	}
	return nil
}
