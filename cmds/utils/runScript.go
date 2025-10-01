package utils

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

var interpreters = map[string][]string{
	"sh":     {"sh", "-c"},
	"bash":   {"bash", "-c"},
	"python": {"python3", "-c"},
	"lua":    {"lua", "-e"},
	"js":     {"node", "-e"},
	"ts":     {"node", "-e"},
}

func RunScript(lang, script string) (stdout, stderr string, err error) {
	args, ok := interpreters[lang]
	if !ok {
		return "", "", fmt.Errorf("unsupported language: %s", lang)
	}
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(args[0], append(args[1:], script)...)
	cmd.Stdout, cmd.Stderr = &outBuf, &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func RunScriptWithTimeout(lang, script string, timeout time.Duration) (stdout, stderr string, err error) {
	args, ok := interpreters[lang]
	if !ok {
		return "", "", fmt.Errorf("unsupported language: %s", lang)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, args[0], append(args[1:], script)...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outBuf, &errBuf
	err = cmd.Run()
	stdout, stderr = outBuf.String(), errBuf.String()
	if ctxErr := cmd.ProcessState.ExitCode(); ctxErr != 0 && err == nil { // none zero exit code is also an error
		err = fmt.Errorf("process exited with code %d: %s", ctxErr, stderr)
	}
	return stdout, stderr, err
}
