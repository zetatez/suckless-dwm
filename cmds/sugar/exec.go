package sugar

import (
	"bytes"
	"os/exec"
)

type ExecService struct{}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) RunScriptShell(script string) (
	stdout string,
	stderr string,
	err error,
) {
	return s.run(
		exec.Command(
			"/bin/bash",
			"-c",
			script,
		),
	)
}

func (s *ExecService) RunScriptPython(script string) (
	stdout string,
	stderr string,
	err error,
) {
	return s.run(
		exec.Command(
			"python3",
			"-c",
			script,
		),
	)
}

func (s *ExecService) RunScriptLua(script string) (
	stdout string,
	stderr string,
	err error,
) {
	return s.run(
		exec.Command(
			"lua",
			"-e",
			script,
		),
	)
}

func (s *ExecService) run(cmd *exec.Cmd) (
	stdout string,
	stderr string,
	err error,
) {
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
