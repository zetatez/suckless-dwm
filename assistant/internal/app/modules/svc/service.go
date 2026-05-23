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

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		home, _ := os.UserHomeDir()
		return path.Join(home, p[2:])
	}
	return p
}

func initClipboard() error {
	var err error
	clipboardOnce.Do(func() {
		err = clipboard.Init()
	})
	return err
}

func (s *Service) notify(msg string) {
	_ = exec.Command("notify-send", msg).Run()
}

func (s *Service) readClipboard() (string, error) {
	if err := initClipboard(); err != nil {
		return "", fmt.Errorf("init clipboard: %w", err)
	}
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}

func (s *Service) writeClipboard(content string) {
	if err := initClipboard(); err != nil {
		s.logger.Warnf("init clipboard: %v", err)
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(content))
}

func (s *Service) isRunning(proc string) bool {
	return exec.Command("pgrep", "-f", proc).Run() == nil
}

func (s *Service) killProcess(proc string) {
	_ = exec.Command("pkill", "-f", proc).Run()
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

var toggleCommands = map[string]string{
	"flameshot":                 "flameshot gui",
	"screenkey":                 "screenkey --key-mode keysyms --opacity 0 -s small --font-color yellow",
	"clipmenu":                  "sh -c clipmenu",
	"passmenu":                  "passmenu",
	"netease-cloud-music":       "netease-cloud-music",
	"htop":                      "st -e htop",
	"yazi":                      "st -e yazi",
	"julia":                     "st -t scratchpad -c scratchpad -e julia",
	"lazydocker":                "st -e lazydocker",
	"lazygit":                   "st -e lazygit",
	"ncmpcpp":                   "st -e ncmpcpp",
	"cava":                      "st -e cava",
	"mutt":                      "st -e mutt",
	"irssi":                     "st -e irssi",
	"newsboat":                  "st -e newsboat",
	"python":                    "st -t scratchpad -c scratchpad -e python -i -c 'import os, sys, datetime, re, json, collections, random, math, numpy as np, pandas as pd, scipy, matplotlib.pyplot as plt; print(dir())'",
	"tty-clock":                 "st -g $GEO -t float -c float -e tty-clock -s",
	"calendar":                  "st -g $GEO -t float -c float -e nvim +':Calendar -view=month'",
	"calendar-scheduling-today": "st -g $GEO -t float -c float -e nvim +':Calendar -view=day'",
}

func (s *Service) Toggle(proc string) string {
	switch {
	case proc == "music":
		return s.toggleMusic()
	case strings.HasPrefix(proc, "rec-"):
		return s.toggleRecording(proc)
	}

	cmd := proc
	if c, ok := toggleCommands[proc]; ok {
		cmd = c
	}
	if strings.Contains(cmd, "$GEO") {
		switch proc {
		case "tty-clock":
			cmd = strings.ReplaceAll(cmd, "$GEO", s.geoForTerminal(0.72, 0.04, 53, 8))
		case "calendar":
			cmd = strings.ReplaceAll(cmd, "$GEO", s.geoForTerminal(0.84, 0.04, 24, 12))
		case "calendar-scheduling-today":
			cmd = strings.ReplaceAll(cmd, "$GEO", s.geoForTerminal(0.80, 0.05, 36, 32))
		}
	}
	if strings.HasPrefix(cmd, "st ") {
		cmd = psl.GetConfig().Svc.DefaultTerminal + cmd[2:]
	}

	match := cmd
	switch proc {
	case "tty-clock", "calendar", "calendar-scheduling-today":
		if i := strings.Index(cmd, " -e "); i >= 0 {
			match = strings.TrimSpace(cmd[i+4:])
		}
	}
	if s.isRunning(match) {
		s.killProcess(match)
		return "killed"
	}
	if proc == "flameshot" {
		_, _, _ = runScript("bash", "systemctl --user start xdg-desktop-portal xdg-desktop-portal-gtk 2>/dev/null || true")
	}
	_ = startScript("bash", cmd)
	return "launched"
}

func (s *Service) toggleMusic() string {
	if s.isRunning("ncmpcpp") || s.isRunning("cava") {
		s.killProcess("ncmpcpp")
		s.killProcess("cava")
		return "killed"
	}
	term := psl.GetConfig().Svc.DefaultTerminal
	_ = startScript("bash", fmt.Sprintf("%s -e ncmpcpp", term))
	_ = startScript("bash", fmt.Sprintf("%s -e cava", term))
	return "launched"
}

func (s *Service) toggleRecording(proc string) string {
	if proc == "rec-show" {
		if s.isRunning("ffplay") {
			s.killProcess("ffplay")
			return "killed"
		}
		term := psl.GetConfig().Svc.DefaultTerminal
		_ = startScript("bash", fmt.Sprintf("%s -e ffplay -loglevel quiet -framedrop -fast -alwaysontop -i /dev/video0", term))
		return "launched"
	}

	if s.isRunning("ffmpeg") {
		s.killProcess("ffmpeg")
		return "killed"
	}

	homeDir, _ := os.UserHomeDir()
	now := time.Now().Format("2006-01-02-15-04-05")
	term := psl.GetConfig().Svc.DefaultTerminal
	scratch := fmt.Sprintf("%s -t scratchpad -c scratchpad -e", term)
	filename := fmt.Sprintf("Videos/rec-%s-%s", strings.TrimPrefix(proc, "rec-"), now)

	switch proc {
	case "rec-audio":
		_ = startScript("bash", fmt.Sprintf("%s ffmpeg -y -r 60 -f alsa -i default -c:a flac %s", scratch, path.Join(homeDir, filename+".flac")))
	case "rec-screen":
		w, h := s.screenSize()
		_ = startScript("bash", fmt.Sprintf("%s ffmpeg -y -s '%dx%d' -r 60 -f x11grab -i %s -f alsa -i default -c:v libx264rgb -crf 0 -preset ultrafast -color_range 2 -c:a aac %s", scratch, w, h, os.Getenv("DISPLAY"), path.Join(homeDir, filename+".mkv")))
	case "rec-webcam":
		_ = startScript("bash", fmt.Sprintf("%s ffmpeg -f pulse -ac 2 -i default -f v4l2 -i /dev/video0 -t 00:00:20 -vcodec libx264 %s", scratch, path.Join(homeDir, filename+".mp4")))
	}
	return "launched"
}

func (s *Service) Launch(command string) error {
	return startScript("bash", command)
}
