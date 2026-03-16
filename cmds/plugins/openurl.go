package plugins

import (
	"fmt"

	"cmds/utils"
)

// chrome --proxy-server="socks5://127.0.0.1:7891" --new-tab "www.chatgpt.com"
func OpenUrlWithChrome(url string) func() {
	return func() {
		_, _, _ = utils.RunScript("bash", fmt.Sprintf(`chrome --proxy-server="%s" "%s"`, ProxyServer, url))
	}
}

// qutebrowser --set content.proxy "socks5://127.0.0.1:7891" "www.chatgpt.com"
func OpenUrlWithQutebrowser(url string) func() {
	return func() {
		_, _, _ = utils.RunScript("bash", fmt.Sprintf(`qutebrowser --set content.proxy "%s" "%s"`, ProxyServer, url))
	}
}

// open url as app
func OpenUrlAsApp(url string) func() {
	return func() {
		// _, _, _ = utils.RunScript("bash", fmt.Sprintf(`chrome --proxy-server="%s" --app="%s"`, ProxyServer, url))
		_, _, _ = utils.RunScript("bash", fmt.Sprintf(`qutebrowser --set content.proxy "%s" --target window "%s"`, ProxyServer, url))
	}
}
