package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var ScriptInterpreters = map[string][]string{
	"sh":     {"sh", "-c"},
	"bash":   {"bash", "-c"},
	"python": {"python3", "-c"},
	"lua":    {"lua", "-e"},
	"js":     {"node", "-e"},
	"ts":     {"node", "-e"},
}

func StartScript(lang, script string) error {
	args, ok := ScriptInterpreters[lang]
	if !ok {
		return fmt.Errorf("unsupported language: %s", lang)
	}

	cmd := exec.Command(args[0], append(args[1:], script)...)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if os.Getenv("DISPLAY") == "" {
		cmd.Env = append(cmd.Env, "DISPLAY=:0")
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

func RunScript(lang, script string) (stdout, stderr string, err error) {
	args, ok := ScriptInterpreters[lang]
	if !ok {
		return "", "", fmt.Errorf("unsupported language: %s", lang)
	}

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(args[0], append(args[1:], script)...)
	cmd.Stdout, cmd.Stderr = &outBuf, &errBuf
	cmd.Env = os.Environ()
	if os.Getenv("DISPLAY") == "" {
		cmd.Env = append(cmd.Env, "DISPLAY=:0")
	}
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func RunScriptWithTimeout(lang, script string, timeout time.Duration) (stdout, stderr string, err error) {
	args, ok := ScriptInterpreters[lang]
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
