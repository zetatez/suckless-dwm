package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Notify(msg ...any) {
	if err := NotifyE(msg...); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "notify failed:", err)
	}
}

func NotifyE(msg ...any) error {
	return NotifyWithOptions(context.Background(), NotifyOptions{
		Title:   defaultNotifyTitle(),
		Message: fmt.Sprint(msg...),
	})
}

type NotifyOptions struct {
	Title   string
	Message string

	// Linux-only.
	AppName string
	Urgency string // "low", "normal", "critical"
	Expire  time.Duration

	// Command execution timeout.
	Timeout time.Duration
}

func NotifyTitleE(title string, msg ...any) error {
	return NotifyWithOptions(context.Background(), NotifyOptions{
		Title:   title,
		Message: fmt.Sprint(msg...),
	})
}

func NotifyWithOptions(ctx context.Context, opt NotifyOptions) error {
	if strings.TrimSpace(opt.Message) == "" {
		return fmt.Errorf("notify: empty message")
	}
	if strings.TrimSpace(opt.Title) == "" {
		opt.Title = defaultNotifyTitle()
	}
	if opt.Timeout <= 0 {
		opt.Timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, opt.Timeout)
	defer cancel()

	switch runtime.GOOS {
	case OSLinux:
		return notifyLinux(ctx, opt)
	case OSDarwin:
		return notifyDarwin(ctx, opt)
	default:
		return fmt.Errorf("notify: unsupported OS: %s", runtime.GOOS)
	}
}

func notifyLinux(ctx context.Context, opt NotifyOptions) error {
	args := make([]string, 0, 10)
	if strings.TrimSpace(opt.AppName) != "" {
		args = append(args, "--app-name", opt.AppName)
	}
	if strings.TrimSpace(opt.Urgency) != "" {
		switch opt.Urgency {
		case "low", "normal", "critical":
			args = append(args, "--urgency", opt.Urgency)
		default:
			return fmt.Errorf("notify: invalid urgency: %q", opt.Urgency)
		}
	}
	if opt.Expire > 0 {
		args = append(args, "--expire-time", fmt.Sprintf("%d", opt.Expire.Milliseconds()))
	}
	args = append(args, opt.Title, opt.Message)

	stderr, err := runCommand(ctx, "notify-send", args...)
	if err != nil {
		if s := strings.TrimSpace(stderr); s != "" {
			return fmt.Errorf("notify-send failed: %w: %s", err, s)
		}
		return fmt.Errorf("notify-send failed: %w", err)
	}
	return nil
}

func notifyDarwin(ctx context.Context, opt NotifyOptions) error {
	script := fmt.Sprintf(
		`display notification "%s" with title "%s"`,
		escapeAppleScriptString(opt.Message),
		escapeAppleScriptString(opt.Title),
	)

	stderr, err := runCommand(ctx, "osascript", "-e", script)
	if err != nil {
		if s := strings.TrimSpace(stderr); s != "" {
			return fmt.Errorf("osascript notify failed: %w: %s", err, s)
		}
		return fmt.Errorf("osascript notify failed: %w", err)
	}
	return nil
}

func runCommand(ctx context.Context, name string, args ...string) (stderr string, err error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	cmd.Stdout = nil

	err = cmd.Run()
	stderr = errBuf.String()
	if ctx.Err() != nil {
		return stderr, fmt.Errorf("command %s cancelled: %w", name, ctx.Err())
	}
	return stderr, err
}

func defaultNotifyTitle() string {
	if len(os.Args) == 0 {
		return "msg"
	}
	base := filepath.Base(os.Args[0])
	base = strings.TrimSpace(base)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "msg"
	}
	return base
}

func escapeAppleScriptString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
