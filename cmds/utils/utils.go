package utils

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/shirou/gopsutil/process"
)

func Lazy(option string, filepath string) {
	switch option {
	case "view", "open", "exec", "copy", "rename", "delete":
		RunScript("bash", fmt.Sprintf("%s -e lazy -o %s -f %s &", GetOSDefaultShell(), option, filepath))
	default:
		return
	}
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
		RunScript("bash", proc)
	}
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

func GetGeoForTerminal(xr float64, yr float64, w int, h int) (geo string) {
	x, y := GetPosition(xr, yr)
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}

func GetGeoCenterForSt(wr float64, hr float64) (geo string) {
	width, height := GetScreenSize()
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
	cmd := fmt.Sprintf("%s -e %s -c '%s'", GetOSDefaultShell(), shell, fmt.Sprintf(`sshpass -p "%s" ssh -o "StrictHostKeyChecking no" -p %d %s@%s`, password, port, user, host))
	_, _, err = RunScript("bash", cmd)
	if err != nil {
		Notify(err)
		return err
	}
	return nil
}
