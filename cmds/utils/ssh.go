package utils

import (
	"fmt"
	"os"
)

func SSH(host string, port int, user string, password string) (err error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	cmd := fmt.Sprintf("%s -e %s -c '%s'", GetOSDefaultTerminal(), shell, fmt.Sprintf(`sshpass -p "%s" ssh -o "StrictHostKeyChecking no" -p %d %s@%s`, password, port, user, host))
	_, _, err = RunScript("bash", cmd)
	if err != nil {
		Notify(err)
		return err
	}
	return nil
}
