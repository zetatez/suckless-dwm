package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var scriptInterpreters = map[string][]string{
	"sh":     {"sh", "-c"},
	"bash":   {"bash", "-c"},
	"python": {"python3", "-c"},
	"lua":    {"lua", "-e"},
	"js":     {"node", "-e"},
	"ts":     {"node", "-e"},
}

func RunScript(lang, script string) (stdout, stderr string, err error) {
	args, ok := scriptInterpreters[lang]
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
	args, ok := scriptInterpreters[lang]
	if !ok {
		return "", "", fmt.Errorf("unsupported language: %s", lang)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var outBuf, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, args[0], append(args[1:], script)...)
	cmd.Stdout, cmd.Stderr = &outBuf, &errBuf
	err = cmd.Run()
	stdout, stderr = outBuf.String(), errBuf.String()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		if err == nil {
			err = context.DeadlineExceeded
		}
		err = fmt.Errorf("script timed out after %s: %w", timeout, err)
	}

	return stdout, stderr, err
}
