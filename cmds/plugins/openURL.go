package plugins

import (
	"fmt"

	"cmds/utils"
)

// chrome --proxy-server=socks5://127.0.0.1:7891  www.chatgpt.com
func ChromeOpenUrl(params, url string) func() {
	return func() {
		utils.RunScript(
			"bash",
			fmt.Sprintf("chrome %s %s", params, url),
		)
	}
}

// edge --kiosk --force-device-scale-factor=1.35 --proxy-server=socks5://127.0.0.1:7891  www.chatgpt.com
func EdgeOpenUrl(params, url string) func() {
	return func() {
		utils.RunScript(
			"bash",
			fmt.Sprintf("microsoft-edge-stable --kiosk --force-device-scale-factor=1.35 %s %s", params, url),
		)
	}
}

func ChromeOpenUrlGoogle() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "http://www.google.com/")()
}

func ChromeOpenUrlChatGPT() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://chatgpt.com/")()
}

func ChromeOpenUrlDouBao() {
	ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.doubao.com/chat/")()
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
