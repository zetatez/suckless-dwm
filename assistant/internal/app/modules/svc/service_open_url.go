package svc

import (
	"fmt"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) OpenURL(browser, url string) error {
	proxy := psl.GetConfig().Svc.ProxyServer
	switch browser {
	case "chrome":
		return startScript("bash", fmt.Sprintf("chrome --proxy-server=%s '%s'", proxy, url))
	case "qutebrowser":
		return startScript("bash", fmt.Sprintf("qutebrowser --set content.proxy '%s' '%s'", proxy, url))
	default:
		return startScript("bash", fmt.Sprintf("xdg-open '%s'", url))
	}
}

func (s *Service) OpenURLAsApp(browser, url string) error {
	proxy := psl.GetConfig().Svc.ProxyServer
	switch browser {
	case "chrome":
		return startScript("bash", fmt.Sprintf("chrome --proxy-server='%s' --app='%s'", proxy, url))
	case "qutebrowser":
		return startScript("bash", fmt.Sprintf("qutebrowser --set content.proxy '%s' --target window '%s'", proxy, url))
	default:
		return startScript("bash", fmt.Sprintf("chrome --proxy-server='%s' --app='%s'", proxy, url))
	}
}
