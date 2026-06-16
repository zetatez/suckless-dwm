package psl

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func StartBackgroundTasks(ctx context.Context) {
	cfg := GetConfig()
	if cfg.Background.Enabled {
		startDaemon(ctx, cfg.Background.Procs)
	}
	startWallpaper(ctx)
}

func startDaemon(ctx context.Context, procs []BackgroundProc) {
	go func() {
		for {
			for _, p := range procs {
				if !isRunning(p.Name) {
					if p.Precursor != "" && !isRunning(p.Precursor) {
						continue
					}
					cmd := exec.CommandContext(ctx, "bash", "-c", p.Command+" &")
					if err := cmd.Start(); err != nil {
						GetLogger().WithError(err).WithField("proc", p.Name).Warn("start proc failed")
						continue
					}
					go func() { _ = cmd.Wait() }()
				}
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}
		}
	}()
}

func startWallpaper(ctx context.Context) {
	go func() {
		dir := GetConfig().Svc.DirWallpaper
		for {
			files, _ := filepath.Glob(filepath.Join(dir, "*.JPG"))
			for _, pic := range files {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("feh --bg-fill '%s'", pic))
				_ = cmd.Run()
				select {
				case <-ctx.Done():
					return
				case <-time.After(60 * time.Second):
				}
			}
			// Empty dir (or dir missing) — wait before scanning again.
			select {
			case <-ctx.Done():
				return
			case <-time.After(60 * time.Second):
			}
		}
	}()
}

func isRunning(name string) bool {
	return exec.Command("pgrep", "-x", name).Run() == nil
}
