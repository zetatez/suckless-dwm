package plugins

import (
	"fmt"

	"cmds/utils"
)

// chrome --proxy-server=socks5://127.0.0.1:7891 www.chatgpt.com
func OpenUrlWithChrome(url string) func() {
	return func() {
		_, _, _ = utils.RunScript("bash", fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 %s", url))
	}
}

// qutebrowser --set content.proxy 'socks5://127.0.0.1:7891' www.chatgpt.com
func OpenUrlWithQutebrowser(url string) func() {
	return func() {
		_, _, _ = utils.RunScript("bash", fmt.Sprintf(`qutebrowser --set content.proxy 'socks5://127.0.0.1:7891' '%s'`, url))
	}
}

// open url as app
func OpenUrlAsApp(url string) func() {
	return func() {
		// _, _, _ = utils.RunScript("bash", fmt.Sprintf("chrome --proxy-server=socks5://127.0.0.1:7891 --app=%s", url))
		_, _, _ = utils.RunScript("bash", fmt.Sprintf(`qutebrowser --set content.proxy 'socks5://127.0.0.1:7891' --target window '%s'`, url))
	}
}
