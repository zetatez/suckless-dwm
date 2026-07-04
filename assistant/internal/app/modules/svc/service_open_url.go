package svc

import (
	"assistant/pkg/utils"
	"fmt"
	"strings"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) OpenURL(browser, url string) error {
	proxy := psl.GetConfig().Settings.VPN
	switch browser {
	case "chrome":
		return utils.StartScript("bash", fmt.Sprintf("chrome --proxy-server=%s '%s'", proxy, url))
	case "qutebrowser":
		return utils.StartScript("bash", fmt.Sprintf("qutebrowser --set content.proxy '%s' '%s'", proxy, url))
	default:
		return utils.StartScript("bash", fmt.Sprintf("xdg-open '%s'", url))
	}
}

func (s *Service) OpenURLAsApp(browser, url string) error {
	proxy := psl.GetConfig().Settings.VPN
	switch browser {
	case "chrome":
		return utils.StartScript("bash", fmt.Sprintf("chrome --proxy-server='%s' --app='%s'", proxy, url))
	case "qutebrowser":
		return utils.StartScript("bash", fmt.Sprintf("qutebrowser --set content.proxy '%s' --target window '%s'", proxy, url))
	default:
		return utils.StartScript("bash", fmt.Sprintf("chrome --proxy-server='%s' --app='%s'", proxy, url))
	}
}

func (s *Service) SearchBooksOnline(query string) error {
	q := strings.ReplaceAll(query, " ", "+")
	urls := []string{
		"https://openlibrary.org/search?q=" + q,
		"https://z-lib.id/s?q=" + q,
	}
	for _, u := range urls {
		if err := s.OpenURL("chrome", u); err != nil {
			s.logger.WithError(err).WithField("url", u).Warn("open books url failed")
		}
	}
	return nil
}

func (s *Service) SearchVideosOnline(query string) error {
	q := strings.ReplaceAll(query, " ", "+")
	urls := []string{
		"https://search.bilibili.com/all?keyword=" + q,
		"https://www.youtube.com/results?search_query=" + q,
	}
	for _, u := range urls {
		if err := s.OpenURL("chrome", u); err != nil {
			s.logger.WithError(err).WithField("url", u).Warn("open videos url failed")
		}
	}
	return nil
}
