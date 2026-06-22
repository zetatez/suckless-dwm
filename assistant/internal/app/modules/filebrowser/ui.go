package filebrowser

import (
	_ "embed"
	"net/http"
)

//go:embed ui.html
var uiHTML string

func renderUI(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(uiHTML))
}
