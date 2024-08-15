package plugins

import (
	"fmt"

	"cmds/sugar"
)

func OpenWeb(params, url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf("chrome %s %s", params, url),
		)
	}
}

func OpenWebGoogle() {
	OpenWeb("--proxy-server="+ProxyServer, "http://www.google.com/")()
}

func OpenWebBing() {
	OpenWeb("", "https://www.bing.com/")()
}

func OpenWebChatGPT() {
	OpenWeb("--proxy-server="+ProxyServer, "https://chatgpt.com/")()
}

func OpenWebDouBao() {
	OpenWeb("", "https://www.doubao.com/chat/")()
}

func OpenWebCodeium() {
	OpenWeb("--proxy-server="+ProxyServer, "https://codeium.com/live/general")()
}

func OpenWebGoogleMail() {
	OpenWeb("--proxy-server="+ProxyServer, "https://mail.google.com/mail")()
}

func OpenWebGoogleTranslate() {
	OpenWeb("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=zh-CN")()
}

func OpenWebGithub() {
	OpenWeb("--proxy-server="+ProxyServer, "https://github.com/zetatez")()
}

func OpenWebGithubGistShareCode() {
	OpenWeb("--proxy-server="+ProxyServer, "https://gist.github.com/")()
}

func OpenWebLeetCode() {
	OpenWeb("", "https://leetcode.cn/search/?q=%E6%9C%80")()
}

func OpenWebWeChat() {
	OpenWeb("", "https://web.wechat.com/")()
}

func OpenWebYouTube() {
	OpenWeb("--proxy-server="+ProxyServer, "https://www.youtube.com")()
}

func OpenWebInstagram() {
	OpenWeb("--proxy-server="+ProxyServer, "https://www.instagram.com")()
}
