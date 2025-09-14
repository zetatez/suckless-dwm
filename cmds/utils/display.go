package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

func GetKeyBoardStatus(kbPath string) (brightness int64, err error) {
	brightnessStr, err := os.ReadFile(kbPath)
	if err != nil {
		return 0, err
	}
	brightness, err = strconv.ParseInt(strings.TrimSpace(string(brightnessStr)), 10, 64)
	if err != nil {
		return 0, err
	}
	return brightness, nil
}

func GetScreenSize() (width int, height int) {
	return robotgo.GetScreenSize()
}

func GetPosition(xr float64, yr float64) (x, y int) {
	width, height := GetScreenSize()
	return int(float64(width) * xr), int(float64(height) * yr)
}

func GetGeoForTerminal(xr float64, yr float64, w int, h int) (geo string) {
	x, y := GetPosition(xr, yr)
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}

func GetGeoCenterForSt(wr float64, hr float64) (geo string) {
	width, height := GetScreenSize()
	w := int(float64(width) * wr)
	h := int(float64(height) * hr)
	x := (width - w) / 2
	y := (height - h) / 2
	return fmt.Sprintf("%dx%d+%d+%d", w, h, x, y)
}
