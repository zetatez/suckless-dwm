package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExecResult struct {
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Success  bool   `json:"success"`
}

type ExecCmd struct {
	command string
	args    []string
	shell   string
}

func ExecCommand(command string, args ...string) *ExecCmd {
	return &ExecCmd{command: command, args: args}
}

func ExecShell(shellCommand string) *ExecCmd {
	return &ExecCmd{shell: shellCommand}
}

func (e *ExecCmd) WithShell(shell string) *ExecCmd {
	e.shell = shell
	return e
}

func (e *ExecCmd) Exec() (*ExecResult, error) {
	var cmd *exec.Cmd
	var display string

	if e.shell != "" {
		shellPath, shellArgs := pickShell(e.shell)
		cmd = exec.Command(shellPath, shellArgs...)
		display = e.shell
	} else {
		cmd = exec.Command(e.command, e.args...)
		display = strings.TrimSpace(fmt.Sprintf("%s %s", e.command, strings.Join(e.args, " ")))
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := &ExecResult{
		Command:  display,
		ExitCode: 0,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Success:  true,
	}

	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		result.Success = false
		return result, nil
	}

	return nil, fmt.Errorf("execute failed: %w", err)
}

func pickShell(shellCommand string) (shellPath string, shellArgs []string) {
	if _, err := os.Stat("/bin/bash"); err == nil {
		return "/bin/bash", []string{"-c", shellCommand}
	}
	if _, err := os.Stat("/bin/sh"); err == nil {
		return "/bin/sh", []string{"-c", shellCommand}
	}
	return shellCommand, nil
}
