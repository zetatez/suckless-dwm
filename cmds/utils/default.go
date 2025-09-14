package utils

func GetOSDefault(objType string) string {
	ObjTypeMap := map[string]func() string{
		"shell":    GetOSDefaultShell,
		"terminal": GetOSDefaultTerminal,
		"editor":   GetOSDefaultEditor,
		"browser":  GetOSDefaultBrowser,
	}
	if fn, ok := ObjTypeMap[objType]; ok {
		return fn()
	}
	Notify("Unsupported object type")
	return ""
}

func GetOSDefaultShell() string {
	OsDefaultShell := map[string]string{
		OsMap["linux"]: "sh",
		OsMap["macos"]: "sh",
	}
	if shell, ok := OsDefaultShell[GetOSType()]; ok {
		return shell
	}
	Notify("Unsupported OS")
	return ""
}

func GetOSDefaultTerminal() string {
	OsDefaultTerminal := map[string]string{
		OsMap["linux"]: "st",
		OsMap["macos"]: "kitty",
	}
	if terminal, ok := OsDefaultTerminal[GetOSType()]; ok {
		return terminal
	}
	Notify("Unsupported OS")
	return ""
}

func GetOSDefaultEditor() string {
	OsDefaultEditor := map[string]string{
		OsMap["linux"]: "nvim",
		OsMap["macos"]: "nvim",
	}
	if editor, ok := OsDefaultEditor[GetOSType()]; ok {
		return editor
	}
	Notify("Unsupported OS")
	return ""
}

func GetOSDefaultBrowser() string {
	OsDefaultBrowser := map[string]string{
		OsMap["linux"]: "chrome",
		OsMap["macos"]: "chrome",
	}
	if browser, ok := OsDefaultBrowser[GetOSType()]; ok {
		return browser
	}
	Notify("Unsupported OS")
	return ""
}
