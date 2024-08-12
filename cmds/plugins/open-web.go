package plugins

import (
	"fmt"

	"cmds/sugar"
)

func ChromeOpenUrl(params, url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf("chrome %s %s", params, url),
		)
	}
}

func EdgeOpenUrl(params, url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf("microsoft-edge-stable %s %s", params, url),
		)
	}
}

// --------------------
func ChromeOpenUrlGoogle() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "http://www.google.com/")()
}

func ChromeOpenUrlBing() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "http://www.bing.com/")()
}

func ChromeOpenUrlChatGPT() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://chatgpt.com/")()
}

func ChromeOpenUrlDouBao() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.doubao.com/chat/")()
}

func ChromeOpenUrlCodeium() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://codeium.com/live/general")()
}

func ChromeOpenUrlGoogleMail() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://mail.google.com/mail")()
}

func ChromeOpenUrlGoogleTranslate() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=zh-CN")()
}

func ChromeOpenUrlGithub() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://github.com/zetatez")()
}

func ChromeOpenUrlGithubGistShareCode() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://gist.github.com/")()
}

func ChromeOpenUrlLeetCode() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://leetcode.cn/search/?q=%E6%9C%80")()
}

func ChromeOpenUrlWeChat() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://web.wechat.com/")()
}

func ChromeOpenUrlYouTube() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.youtube.com")()
}

func ChromeOpenUrlInstagram() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.instagram.com")()
}

// --------------------
func EdgeOpenUrlGoogle() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "http://www.google.com/")()
}

func EdgeOpenUrlBing() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "http://www.bing.com/")()
}

func EdgeOpenUrlChatGPT() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://chatgpt.com/")()
}

func EdgeOpenUrlDouBao() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://www.doubao.com/chat/")()
}

func EdgeOpenUrlCodeium() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://codeium.com/live/general")()
}

func EdgeOpenUrlGoogleMail() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://mail.google.com/mail")()
}

func EdgeOpenUrlGoogleTranslate() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=zh-CN")()
}

func EdgeOpenUrlGithub() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://github.com/zetatez")()
}

func EdgeOpenUrlGithubGistShareCode() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://gist.github.com/")()
}

func EdgeOpenUrlLeetCode() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://leetcode.cn/search/?q=%E6%9C%80")()
}

func EdgeOpenUrlWeChat() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://web.wechat.com/")()
}

func EdgeOpenUrlYouTube() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://www.youtube.com")()
}

func EdgeOpenUrlInstagram() {
	EdgeOpenUrl("--proxy-server="+ProxyServer, "https://www.instagram.com")()
}
