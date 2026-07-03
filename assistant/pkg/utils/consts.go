package utils

const (
	OSLinux   = "linux"
	OSDarwin  = "darwin"
	OSWindows = "windows"
)

var OsMap = map[string]string{
	"linux":   OSLinux,
	"macos":   OSDarwin,
	"windows": OSWindows,
}
