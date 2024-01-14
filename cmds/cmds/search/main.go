package main

import (
	"cmds/plugins"
	"cmds/sugar"
)

func main() {
	NewSearch().Search()
}

var ActionMap = map[string]func(){
	"launchapp: passmenu":     sugar.ReturnLaunchApp("passmenu"),
	"launchapp: chrome":       sugar.ReturnLaunchApp("chrome --proxy-server='socks5://127.0.0.1:7891'"),
	"launchapp: edge":         sugar.ReturnLaunchApp("microsoft-edge-stable"),
	"workflow: format-sql":    plugins.ReturnFormatSql(),
	"workflow: format-json":   plugins.ReturnFormatJson(),
	"website: trans":          sugar.ReturnWebsite("https://translate.google.com/?hl=zh-CN"),
	"website: suckless":       sugar.ReturnWebsite("https://dwm.suckless.org"),
	"website: mirror":         sugar.ReturnWebsite("https://developer.aliyun.com/mirror"),
	"website: arch wiki":      sugar.ReturnWebsite("https://wiki.archlinux.org"),
	"website: arxiv":          sugar.ReturnWebsite("https://arxiv.org"),
	"website: bili":           sugar.ReturnWebsite("https://www.bilibili.com"),
	"website: bing":           sugar.ReturnWebsite("https://cn.bing.com"),
	"website: cctv5":          sugar.ReturnWebsite("https://tv.cctv.com/live/cctv5"),
	"website: github":         sugar.ReturnWebsite("https://github.com/zetatez?tab=repositories"),
	"website: mall":           sugar.ReturnWebsite("https://www.jd.com"),
	"website: map":            sugar.ReturnWebsite("https://www.google.com/maps/place/shanghai"),
	"website: news":           sugar.ReturnWebsite("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"website: ocr":            sugar.ReturnWebsite("http://ocr.space"),
	"website: scholar":        sugar.ReturnWebsite("https://scholar.google.com"),
	"website: wolframalpha":   sugar.ReturnWebsite("https://www.wolframalpha.com"),
	"website: youtube":        sugar.ReturnWebsite("https://www.youtube.com"),
	"website: runoob":         sugar.ReturnWebsite("https://www.runoob.com"),
	"website: ajax":           sugar.ReturnWebsite("https://www.runoob.com/ajax/ajax-tutorial.html"),
	"website: angular":        sugar.ReturnWebsite("https://www.runoob.com/angularjs2/angularjs2-tutorial.html"),
	"website: css":            sugar.ReturnWebsite("https://www.runoob.com/css3/css3-tutorial.html"),
	"website: design pattern": sugar.ReturnWebsite("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"website: docker":         sugar.ReturnWebsite("https://www.runoob.com/docker/docker-tutorial.html"),
	"website: html":           sugar.ReturnWebsite("https://www.runoob.com/html/html5-intro.html"),
	"website: javascript":     sugar.ReturnWebsite("https://www.runoob.com/js/js-tutorial.html"),
	"website: maven":          sugar.ReturnWebsite("https://www.runoob.com/maven/maven-tutorial.html"),
	"website: mongo":          sugar.ReturnWebsite("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"website: nodejs":         sugar.ReturnWebsite("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"website: react":          sugar.ReturnWebsite("https://www.runoob.com/react/react-tutorial.html"),
	"website: redis":          sugar.ReturnWebsite("https://www.runoob.com/redis/redis-tutorial.html"),
	"website: regex":          sugar.ReturnWebsite("https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference"),
	"website: typescript":     sugar.ReturnWebsite("https://www.runoob.com/typescript/ts-tutorial.html"),
	"website: vue":            sugar.ReturnWebsite("https://www.runoob.com/vue3/vue3-tutorial.html"),
}

type Search struct{}

func NewSearch() *Search {
	return &Search{}
}

func (s *Search) Search() {
	list := make([]string, 0)
	for k := range ActionMap {
		list = append(list, k)
	}
	content, err := sugar.Choose(list)
	if err != nil {
		return
	}
	f, ok := ActionMap[content]
	switch {
	case ok:
		f()
		return
	case sugar.Exists(content) && sugar.IsFile(content):
		sugar.Lazy("open", content)
		return
	case sugar.IsUrl(content):
		sugar.ReturnWebsite(content)()
		return
	default:
		sugar.SearchFromWeb(content)
	}
}
