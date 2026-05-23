package psl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	backgroundOnce sync.Once
	wallpaperOnce  sync.Once
)

func StartBackgroundTasks() {
	cfg := GetConfig()
	if cfg.Background.Enabled {
		startDaemon(cfg.Background.Procs)
	}
	startWallpaper()
}

func startDaemon(procs []BackgroundProc) {
	backgroundOnce.Do(func() {
		go func() {
			for {
				for _, p := range procs {
					if !isRunning(p.Name) {
						cmd := exec.Command("bash", "-c", p.Command+" &")
						cmd.Start()
						go cmd.Wait()
					}
				}
				time.Sleep(3 * time.Second)
			}
		}()
	})
}

func startWallpaper() {
	wallpaperOnce.Do(func() {
		go func() {
			dir := GetConfig().Svc.DirWallpaper
			if strings.HasPrefix(dir, "~/") {
				home, _ := os.UserHomeDir()
				dir = filepath.Join(home, dir[2:])
			}
			for {
				files, err := filepath.Glob(filepath.Join(dir, "*.JPG"))
				if err != nil || len(files) == 0 {
					time.Sleep(60 * time.Second)
					continue
				}
				for _, pic := range files {
					_, _, _ = func() (string, string, error) {
						cmd := exec.Command("bash", "-c", fmt.Sprintf("feh --bg-fill '%s'", pic))
						var o, e strings.Builder
						cmd.Stdout, cmd.Stderr = &o, &e
						err := cmd.Run()
						return o.String(), e.String(), err
					}()
					time.Sleep(60 * time.Second)
				}
			}
		}()
	})
}

func isRunning(name string) bool {
	return exec.Command("pgrep", "-x", name).Run() == nil
}
