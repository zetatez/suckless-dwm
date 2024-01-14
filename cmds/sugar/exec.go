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

func (s *ExecService) RunScript(
	scriptType string,
	script string,
) (stdout string, stderr string, err error) {
	switch scriptType {
	case "shell":
		return s.RunScriptShell(script)
	case "python":
		return s.RunScriptPython(script)
	case "lua":
		return s.RunScriptLua(script)
	default:
		err = fmt.Errorf("script type %v not supported", scriptType)
		return "", "", err
	}
}

func (s *ExecService) RunScriptShell(script string) (stdout string, stderr string, err error) {
	cmd := exec.Command("/bin/bash", "-c", script)
	return s.exec(cmd)
}

func (s *ExecService) RunScriptPython(script string) (stdout string, stderr string, err error) {
	cmd := exec.Command("python3", "-c", script)
	return s.exec(cmd)
}

func (s *ExecService) RunScriptLua(script string) (stdout string, stderr string, err error) {
	cmd := exec.Command("lua", "-e", script)
	return s.exec(cmd)
}

func (s *ExecService) exec(cmd *exec.Cmd) (stdout string, stderr string, err error) {
	var stdoutbyte, stderrbyte bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdoutbyte, &stderrbyte
	err = cmd.Run()
	if err != nil {
		stdout, stderr = stdoutbyte.String(), stderrbyte.String()
		return stdout, stderr, err
	}
	stdout, stderr = stdoutbyte.String(), stderrbyte.String()
	return stdout, stderr, nil
}
