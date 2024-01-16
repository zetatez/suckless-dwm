package plugins

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"cmds/sugar"
)

func SshTo() {
	mysshListFilePath := path.Join(os.Getenv("HOME"), ".ssh/my.ssh.list")
	if !sugar.IsFileExists(mysshListFilePath) {
		f, err := os.Create(mysshListFilePath)
		if err != nil {
			sugar.Notify(err)
			return
		}
		f.Close()
	}

	// read from to ~/.ssh/my.ssh.list
	b, err := os.ReadFile(mysshListFilePath)
	if err != nil {
		sugar.Notify(err)
		return
	}
	mySshList := []map[string]string{}
	slice1 := strings.Split(string(b), "\n")
	for _, x := range slice1 {
		x = strings.TrimSpace(x)
		slice2 := regexp.MustCompile(`[ \r\t\s]+`).Split(x, -1)
		if len(slice2) < 3 {
			continue
		}
		host := strings.TrimSpace(slice2[0])
		user := strings.TrimSpace(slice2[1])
		password := strings.TrimSpace(slice2[2])
		slice3 := strings.Split(x, "#")
		if len(slice3) != 2 {
			continue
		}
		desc := strings.TrimSpace(slice3[1])
		mySshList = append(
			mySshList,
			map[string]string{"host": host, "user": user, "password": password, "desc": desc},
		)
	}

	// read from ~/.ssh/known_hosts
	knownHosts, err := sugar.GetKnownHosts()
	if err != nil {
		sugar.Notify(err)
		return
	}

	// prompt
	promptList := []string{}
	for _, x := range mySshList {
		promptList = append(promptList, fmt.Sprintf("%-20s %-20s %-20s # %s", x["host"], x["user"], x["password"], x["desc"]))
	}
	for host := range knownHosts {
		promptList = append(promptList, fmt.Sprintf("%-20s", knownHosts[host]))
	}

	// choose
	chioce, err := sugar.Choose("ssh to: ", promptList)
	if err != nil {
		sugar.Notify(err)
		return
	}
	chioce = strings.TrimSpace(chioce)
	slice := regexp.MustCompile(`[ \r\t\s]+`).Split(chioce, -1)

	switch {
	case len(slice) > 3:
		host := strings.TrimSpace(slice[0])
		user := strings.TrimSpace(slice[1])
		password := strings.TrimSpace(slice[2])
		err = sugar.SSH(host, 22, user, password)
		if err != nil {
			sugar.Notify(err)
			return
		}
		return
	default:
		host := strings.TrimSpace(slice[0])
		user, err := sugar.GetInput("user: ")
		if err != nil {
			sugar.Notify(err)
			return
		}
		password, err := sugar.GetInput("password: ")
		if err != nil {
			sugar.Notify(err)
			return
		}
		desc, err := sugar.GetInput("desc: ")
		if err != nil {
			sugar.Notify(err)
			return
		}

		err = sugar.SSH(host, 22, user, password)
		if err != nil {
			sugar.Notify(err)
			return
		}

		// append to ~/.ssh/my.ssh.list
		file, err := os.OpenFile(
			mysshListFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o666,
		)
		if err != nil {
			sugar.Notify(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		fmt.Fprintf(writer, "%-20s %-20s %-20s # %s\r\n", host, user, password, desc)
		writer.Flush()
	}
}
