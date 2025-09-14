package utils

import (
	"bytes"
	"fmt"
	"os/exec"
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
