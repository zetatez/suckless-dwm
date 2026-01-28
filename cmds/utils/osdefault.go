package utils

var OsDefaultMap = map[string]map[string]string{
	"shell": {
		"linux":  "sh",
		"darwin": "sh",
	},
	"terminal": {
		"linux":  "st",
		"darwin": "kitty",
	},
	"editor": {
		"linux":  "nvim",
		"darwin": "nvim",
	},
	"browser": {
		"linux":  "chrome",
		"darwin": "chrome",
	},
}

func GetOSDefault(objType string) string {
	if m, ok := OsDefaultMap[objType]; ok {
		if v, ok := m[GetOSType()]; ok {
			return v
		}
	}
	Notify("Unsupported: " + objType)
	return ""
}

func GetOSDefaultShell() string {
	return GetOSDefault("shell")
}

func GetOSDefaultTerminal() string {
	return GetOSDefault("terminal")
}

func GetOSDefaultEditor() string {
	return GetOSDefault("editor")
}

func GetOSDefaultBrowser() string {
	return GetOSDefault("browser")
}
