package svc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"assistant/internal/bootstrap/psl"

	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

type Service struct {
	logger *logrus.Logger
}

func NewService() *Service {
	return &Service{logger: psl.GetLogger()}
}

var interpreters = map[string][]string{
	"sh":     {"sh", "-c"},
	"bash":   {"bash", "-c"},
	"python": {"python3", "-c"},
}

func runScript(lang, script string) (string, string, error) {
	args, ok := interpreters[lang]
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
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func startScript(lang, script string) error {
	args, ok := interpreters[lang]
	if !ok {
		return fmt.Errorf("unsupported language: %s", lang)
	}
	cmd := exec.Command(args[0], append(args[1:], script)...)
	cmd.Env = os.Environ()
	if os.Getenv("DISPLAY") == "" {
		cmd.Env = append(cmd.Env, "DISPLAY=:0")
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

var clipboardOnce sync.Once

func initClipboard() error {
	var err error
	clipboardOnce.Do(func() {
		err = clipboard.Init()
	})
	return err
}

func (s *Service) notify(msg string) {
	cmd := exec.Command("notify-send", msg)
	// CombinedOutput absorbs stdout/stderr so the child never blocks on a
	// closed pipe and notify-send's own errors don't propagate.
	_, _ = cmd.CombinedOutput()
}

func (s *Service) readClipboard() (string, error) {
	if err := initClipboard(); err != nil {
		return "", fmt.Errorf("init clipboard: %w", err)
	}
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}

func (s *Service) writeClipboard(content string) error {
	if err := initClipboard(); err != nil {
		return fmt.Errorf("init clipboard: %w", err)
	}
	clipboard.Write(clipboard.FmtText, []byte(content))
	return nil
}

// pushClipboard writes content to the clipboard, surfaces a notification, and
// returns a copy of the value for the HTTP response. Callers that need to
// surface a clipboard failure to the client should inspect the returned
// error; fire-and-forget callers should use copyToClipboardWithNotify.
func (s *Service) pushClipboard(value, summary string) (string, error) {
	if err := s.writeClipboard(value); err != nil {
		s.logger.WithError(err).Warn("write clipboard failed")
		return value, err
	}
	s.notify(summary)
	return value, nil
}

// copyToClipboardWithNotify is the fire-and-forget variant of pushClipboard:
// any clipboard error is logged and swallowed, since these callers return
// their own values regardless of clipboard state.
func (s *Service) copyToClipboardWithNotify(value, summary string) {
	if err := s.writeClipboard(value); err != nil {
		s.logger.WithError(err).Warn("write clipboard failed")
		return
	}
	s.notify(summary)
}

// rofiPrompt opens rofi with the given prompt and returns the trimmed result.
// An empty string is returned if the user dismisses the dialog.
func (s *Service) rofiPrompt(prompt string) (string, error) {
	out, _, err := runScript("bash", fmt.Sprintf("printf '' | rofi -dmenu -p '%s'", prompt))
	return strings.TrimSpace(out), err
}

func (s *Service) isRunning(proc string) bool {
	return exec.Command("pgrep", proc).Run() == nil
}

func (s *Service) killProcess(proc string) error {
	if proc == "" {
		return fmt.Errorf("proc cannot be empty")
	}
	return exec.Command("pkill", proc).Run()
}

func (s *Service) screenSize() (int, int) {
	out, _, err := runScript("bash", "xdpyinfo|awk '/dimensions/{split($2,a,\"x\");print a[1],a[2]}'")
	if err != nil {
		return 1920, 1080
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) < 2 {
		return 1920, 1080
	}
	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])
	if w <= 0 || h <= 0 {
		return 1920, 1080
	}
	return w, h
}

func (s *Service) geoForTerminal(xr, yr float64, w, h int) string {
	sw, sh := s.screenSize()
	x := int(float64(sw) * xr)
	y := int(float64(sh) * yr)
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}

func (s *Service) Launch(command string) error {
	return startScript("bash", command)
}

func (s *Service) Toggle(cmd, match string) (string, error) {
	if match == "" {
		return "", fmt.Errorf("match cannot be empty")
	}
	s.logger.Info(fmt.Sprintf("toggle %s", cmd))
	if s.isRunning(match) {
		s.logger.Info("running")
		if err := s.killProcess(match); err != nil {
			s.logger.Info("killed")
			return "", fmt.Errorf("kill %s: %w", match, err)
		}
		return "killed", nil
	}
	if err := startScript("bash", cmd); err != nil {
		return "", fmt.Errorf("launch %s: %w", cmd, err)
	}
	return "launched", nil
}

func (s *Service) ToggleTTYClock() (string, error) {
	term := psl.GetConfig().Svc.DefaultTerminal
	geo := s.geoForTerminal(0.72, 0.04, 53, 8)
	cmd := fmt.Sprintf("%s -g %s -t float -c float -e tty-clock -s", term, geo)
	return s.Toggle(cmd, "tty-clock")
}

func (s *Service) ToggleMusic() (string, error) {
	if s.isRunning("ncmpcpp") || s.isRunning("cava") {
		if err := s.killProcess("ncmpcpp"); err != nil {
			return "", fmt.Errorf("kill ncmpcpp: %w", err)
		}
		if err := s.killProcess("cava"); err != nil {
			return "", fmt.Errorf("kill cava: %w", err)
		}
		return "killed", nil
	}
	term := psl.GetConfig().Svc.DefaultTerminal
	if err := startScript("bash", fmt.Sprintf("%s -e ncmpcpp", term)); err != nil {
		return "", fmt.Errorf("launch ncmpcpp: %w", err)
	}
	if err := startScript("bash", fmt.Sprintf("%s -e cava", term)); err != nil {
		return "", fmt.Errorf("launch cava: %w", err)
	}
	return "launched", nil
}

func (s *Service) ToggleRecShow() (string, error) {
	term := psl.GetConfig().Svc.DefaultTerminal
	cmd := fmt.Sprintf("%s -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0", term)
	return s.Toggle(cmd, "ffplay")
}

func (s *Service) ToggleRecAudio() (string, error) {
	if s.isRunning("ffmpeg") {
		if err := s.killProcess("ffmpeg"); err != nil {
			return "", fmt.Errorf("kill ffmpeg: %w", err)
		}
		return "killed", nil
	}
	homeDir, _ := os.UserHomeDir()
	now := time.Now().Format("2006-01-02-15-04-05")
	term := psl.GetConfig().Svc.DefaultTerminal
	scratch := fmt.Sprintf("%s -t scratchpad -c scratchpad -e", term)
	filename := path.Join(homeDir, fmt.Sprintf("Videos/rec-audio-%s.flac", now))
	cmd := fmt.Sprintf("%s ffmpeg -y -r 60 -f alsa -i default -c:a flac %s", scratch, filename)
	return s.Toggle(cmd, "ffmpeg")
}

func (s *Service) ToggleRecScreen() (string, error) {
	if s.isRunning("ffmpeg") {
		if err := s.killProcess("ffmpeg"); err != nil {
			return "", fmt.Errorf("kill ffmpeg: %w", err)
		}
		return "killed", nil
	}
	homeDir, _ := os.UserHomeDir()
	now := time.Now().Format("2006-01-02-15-04-05")
	term := psl.GetConfig().Svc.DefaultTerminal
	scratch := fmt.Sprintf("%s -t scratchpad -c scratchpad -e", term)
	filename := path.Join(homeDir, fmt.Sprintf("Videos/rec-screen-%s.mkv", now))
	w, h := s.screenSize()
	cmd := fmt.Sprintf("%s ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s", scratch, w, h, os.Getenv("DISPLAY"), filename)
	return s.Toggle(cmd, "ffmpeg")
}

func (s *Service) ToggleRecWebcam() (string, error) {
	if s.isRunning("ffmpeg") {
		if err := s.killProcess("ffmpeg"); err != nil {
			return "", fmt.Errorf("kill ffmpeg: %w", err)
		}
		return "killed", nil
	}
	homeDir, _ := os.UserHomeDir()
	now := time.Now().Format("2006-01-02-15-04-05")
	term := psl.GetConfig().Svc.DefaultTerminal
	scratch := fmt.Sprintf("%s -t scratchpad -c scratchpad -e", term)
	filename := path.Join(homeDir, fmt.Sprintf("Videos/rec-webcam-%s.mp4", now))
	cmd := fmt.Sprintf("%s ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s", scratch, filename)
	return s.Toggle(cmd, "ffmpeg")
}
