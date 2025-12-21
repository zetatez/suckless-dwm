package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	wallpaperDir := filepath.Join(os.Getenv("HOME"), "Pictures", "wallpapers")

	for {
		files, err := filepath.Glob(filepath.Join(wallpaperDir, "*.JPG"))
		if err != nil {
			fmt.Println("failed to get files:", err)
			return
		}

		if len(files) == 0 {
			time.Sleep(30 * time.Second)
		}

		for _, picture := range files {
			cmd := exec.Command("feh", "--bg-fill", picture)
			if err := cmd.Run(); err != nil {
				fmt.Println("failed to set wallpaper:", err)
			}
			time.Sleep(30 * time.Second)
		}
	}
}
