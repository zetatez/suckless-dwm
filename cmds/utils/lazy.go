package utils

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

func Lazy(option string, filepath string) {
	switch option {
	case "view", "open", "exec", "copy", "rename", "delete":
		RunScript("bash", fmt.Sprintf("%s -e lazy -o %s -f %s &", GetOSDefaultTerminal(), option, filepath))
	default:
		return
	}
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
