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

func OpenWebGoogleTranslate() {
	OpenWeb("https://translate.google.com/?sl=auto&tl=zh-CN")()
}

func OpenWebChatGPT() {
	OpenWeb("https://chatgpt.com/")()
}

func OpenWebGoogleMail() {
	OpenWeb("https://mail.google.com/mail")()
}

func OpenWebLeetCode() {
	OpenWeb("https://leetcode.cn/problemset/")()
}

func OpenWebYouTube() {
	OpenWeb("https://www.youtube.com")()
}

func OpenWebInstagram() {
	OpenWeb("https://www.instagram.com")()
}

func OpenWebGithubGistShareCode() {
	OpenWeb("https://gist.github.com/")()
}
