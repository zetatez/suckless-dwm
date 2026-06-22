package kindle

import (
	_ "embed"
	"html/template"
	"net/http"
	"time"
)

//go:embed ui.html
var uiHTML string

//go:embed content.html
var contentHTML string

type uiData struct {
	Clock    string
	Date     string
	News     template.HTML
	Calendar template.HTML
	Market   template.HTML
}

func renderUI(svc *Service, w http.ResponseWriter) error {
	now := time.Now()
	clock := now.Format("15:04")
	date := now.Format("2006年1月2日  Monday")
	weather := svc.FetchWeather()
	if weather != "" {
		date = date + "  ·  " + weather
	}

	tmpl, err := template.New("ui").Parse(uiHTML)
	if err != nil {
		return err
	}

	data := uiData{
		Clock:    clock,
		Date:     date,
		News:     template.HTML(svc.BuildNewsHTML()),
		Calendar: template.HTML(svc.BuildCalendar(now)),
		Market:   template.HTML(svc.MarketChart()),
	}
	return tmpl.Execute(w, data)
}

func renderContent(svc *Service, w http.ResponseWriter) error {
	now := time.Now()
	clock := now.Format("15:04")
	date := now.Format("2006年1月2日  Monday")
	weather := svc.FetchWeather()
	if weather != "" {
		date = date + "  ·  " + weather
	}

	tmpl, err := template.New("content").Parse(contentHTML)
	if err != nil {
		return err
	}

	data := uiData{
		Clock:    clock,
		Date:     date,
		News:     template.HTML(svc.BuildNewsHTML()),
		Calendar: template.HTML(svc.BuildCalendar(now)),
		Market:   template.HTML(svc.MarketChart()),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}
