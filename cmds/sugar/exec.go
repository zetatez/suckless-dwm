package sugar

import (
	"bytes"
	"fmt"
	"os/exec"
)

type ExecService struct{}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) RunScript(lang, script string) (stdout, stderr string, err error) {
	var cmd *exec.Cmd

	switch lang {
	case "bash":
		cmd = exec.Command("/bin/bash", "-c", script)
	case "python":
		cmd = exec.Command("python3", "-c", script)
	case "lua":
		cmd = exec.Command("lua", "-e", script)
	default:
		err = fmt.Errorf("unsupported language: %s", lang)
		return "", "", err
	}

	return s.run(cmd)
}

func (s *ExecService) run(cmd *exec.Cmd) (stdout, stderr string, err error) {
	var stdoutbyte, stderrbyte bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdoutbyte, &stderrbyte
	err = cmd.Run()
	stdout, stderr = stdoutbyte.String(), stderrbyte.String()
	return stdout, stderr, err
}
