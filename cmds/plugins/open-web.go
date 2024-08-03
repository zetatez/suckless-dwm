package plugins

import (
	"fmt"

	"cmds/sugar"
)

func OpenWeb(url string) func() {
	return func() {
		sugar.NewExecService().RunScriptShell(
			fmt.Sprintf(
				"chrome --proxy-server=%s %s",
				ProxyServer,
				url,
			),
		)
	}
}

func OpenWebChatGPT() {
	OpenWeb("https://chatgpt.com/")()
}

func OpenWebDouBao() {
	OpenWeb("https://www.doubao.com/chat/")()
}

func OpenWebGemini() {
	OpenWeb("https://gemini.google.com/app")()
}

func OpenWebGoogleMail() {
	OpenWeb("https://mail.google.com/mail")()
}

func OpenWebGoogleTranslate() {
	OpenWeb("https://translate.google.com/?sl=auto&tl=zh-CN")()
}

func OpenWebGithub() {
	OpenWeb("https://github.com/zetatez")()
}

func OpenWebGithubGistShareCode() {
	OpenWeb("https://gist.github.com/")()
}

func OpenWebLeetCode() {
	OpenWeb("https://leetcode.cn/search/?q=%E6%9C%80")()
}

func OpenWebWeChat() {
	OpenWeb("https://web.wechat.com/")()
}

func OpenWebYouTube() {
	OpenWeb("https://www.youtube.com")()
}

func OpenWebInstagram() {
	OpenWeb("https://www.instagram.com")()
}
