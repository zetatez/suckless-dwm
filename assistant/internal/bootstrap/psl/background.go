package psl

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"assistant/pkg/dwmblocknotify"
	"assistant/pkg/news"
)

func StartBackgroundTasks(ctx context.Context) {
	cfg := GetConfig()
	if cfg.Background.Enabled {
		startDaemon(ctx, cfg.Background.Procs)
	}
	startWallpaper(ctx)
	// startNewsNotify(ctx)
}

func startDaemon(ctx context.Context, procs []BackgroundProc) {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
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
			case <-ticker.C:
			}
		}
	}()
}

func isRunning(name string) bool {
	return exec.Command("pgrep", "-x", name).Run() == nil
}

func startWallpaper(ctx context.Context) {
	go func() {
		dir := GetConfig().Settings.DirWallpaper
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
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
				case <-ticker.C:
				}
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

func startNewsNotify(ctx context.Context) {
	go func() {
		collector := news.New()
		items, err := collector.Fetch(ctx, "top-news", 20)
		if err != nil {
			GetLogger().WithError(err).Warn("fetch news failed")
		}

		fetchTicker := time.NewTicker(time.Minute * 30)
		defer fetchTicker.Stop()
		sendTicker := time.NewTicker(16 * time.Second)
		defer sendTicker.Stop()
		idx := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-fetchTicker.C:
				newsItems, err := collector.Fetch(ctx, "top-news", 20)
				if err != nil {
					GetLogger().WithError(err).Warn("fetch news failed")
					continue
				}
				items = newsItems
				idx = 0
			case <-sendTicker.C:
				if len(items) == 0 {
					continue
				}
				item := items[idx]
				msg := item.Title
				// if item.Source != "" {
				// 	msg = item.Source + " | " + msg
				// }
				dwmblocknotify.POST(msg, 3*time.Second)
				idx = (idx + 1) % len(items)
			}
		}
	}()
}
