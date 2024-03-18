package plugins

import (
	"sort"

	"cmds/sugar"
)

var ActionMap = map[string]func(){
	"lazy open search file content":   LazyOpenSearchFileContent,
	"lazy open search file":           LazyOpenSearchFile,
	"lazy open search book":           LazyOpenSearchBook,
	"lazy open search wiki":           LazyOpenSearchWiki,
	"lazy open search media":          LazyOpenSearchMedia,
	"ocr":                             OCR,
	"ssh-to":                          SshTo,
	"umount xyz":                      UmountXYZ,
	"handle-copied":                   HandleCopied,
	"wifi-connect":                    WifiConnect,
	"format sql":                      FormatSql,
	"format json":                     FormatJson,
	"format yaml":                     FormatYaml,
	"search books online":             SearchBooksOnline,
	"search videos online":            SearchVideosOnline,
	"get current-datetime":            GetCurrentDatetime,
	"get current-unix-sec":            GetCurrentUnixSec,
	"get host name":                   GetHostName,
	"get ip":                          GetIP,
	"excalidraw":                      OpenWeb("http://127.0.0.1:3000"), // cd ~/github/excalidraw && docker-compose up -d
	"map: google ap":                  OpenWeb("https://www.google.com/maps/place/shanghai"),
	"email: 163":                      OpenWeb("https://mail.163.com"),
	"email: gmail":                    OpenWeb("https://accounts.google.com/b/0/AddMailService"),
	"email: outlook":                  OpenWeb("https://outlook.live.com/mail"),
	"translate auto to en":            OpenWeb("https://translate.google.com/?sl=auto&tl=en"),    // OpenWeb("https://fanyi.baidu.com/#auto/en/"),
	"translate auto to zh":            OpenWeb("https://translate.google.com/?sl=auto&tl=zh-CN"), // OpenWeb("https://fanyi.baidu.com/#auto/zh/"),
	"transform date time to unix sec": TransformDatetime2UnixSec,
	"transform unix sec to date time": TransformUnixSec2DateTime,
	"note: diary":                     NoteDiary,
	"note: timeline":                  NoteTimeline,
	"note: flash card":                NoteFlashCard,
	"launch app: baidudisknet":        LaunchApp("baidudisknet"),
	"launch app: chrome":              LaunchApp("chrome --proxy-server=socks5://127.0.0.1:7891"),
	"launch app: edge":                LaunchApp("edge --proxy-server=socks5://127.0.0.1:7891"),
	"launch app: dingtalk":            LaunchApp("dingtalk"),
	"launch app: inkscape":            LaunchApp("inkscape"),
	"launch app: krita":               LaunchApp("krita"),
	"launch app: netease-cloud-music": LaunchApp("netease-cloud-music"),
	"launch app: obsidian":            LaunchApp("obsidian"),
	"launch app: passmenu":            LaunchApp("passmenu"),
	"launch app: passmmenu":           LaunchApp("passmmenu"),
	"launch app: scribus":             LaunchApp("scribus"),
	"launch app: slack":               LaunchApp("slack"),
	"launch app: subl":                LaunchApp("subl"),
	"launch app: wechat-uos":          LaunchApp("wechat-uos"),
	"launch app: wemeet":              LaunchApp("wemeet"),
	"launch app: wps":                 LaunchApp("wps"),
	"launch app: xournalpp":           LaunchApp("xournalpp"),
	"launch app: zoom":                LaunchApp("zoom"),
	"toggle: address-book":            ToggleAddressbook,
	"toggle: bluetooth":               ToggleBlueTooth,
	"toggle: calendar today schedule": ToggleCalendarTodaySchedule,
	"toggle: calendar scheduling":     ToggleCalendarScheduling,
	"toggle: chrome":                  ToggleChrome,
	"toggle: edge":                    ToggleEdge,
	"toggle: clipmenu":                ToggleClipmenu,
	"toggle: flameshot":               ToggleFlameshot,
	"toggle: inkscape":                ToggleInkscape,
	"toggle: irssi":                   ToggleIrssi,
	"toggle: joshuto":                 ToggleJoshuto,
	"toggle: julia":                   ToggleJulia,
	"toggle: keyboard-light":          ToggleKeyboardLight,
	"toggle: krita":                   ToggleKrita,
	"toggle: lazydocker":              ToggleLazyDocker,
	"toggle: lua":                     ToggleLua,
	"toggle: music":                   ToggleMusic,
	"toggle: music-net-cloud":         ToggleMusicNetCloud,
	"toggle: mutt":                    ToggleMutt,
	"toggle: passmenu":                TogglePassmenu,
	"toggle: python":                  TogglePython,
	"toggle: rec-audio":               ToggleRecAudio,
	"toggle: rec-screen":              ToggleRecScreen,
	"toggle: rec-webcam":              ToggleRecWebcam,
	"toggle: redshift":                ToggleRedShift,
	"toggle: scala":                   ToggleScala,
	"toggle: screen":                  ToggleScreen,
	"toggle: screenkey":               ToggleScreenKey,
	"toggle: show":                    ToggleShow,
	"toggle: sublime":                 ToggleSublime,
	"toggle: sys-short-cuts":          ToggleSysShortcuts,
	"toggle: top":                     ToggleTop,
	"toggle: wallpaper":               ToggleWallpaper,
	"toggle: wechat":                  ToggleWechat,
	"toggle: xournal":                 ToggleXournal,
	"jump to code from log":           JumpToCodeFromLog,
	"www: wechat file helper":         OpenWeb("https://filehelper.weixin.qq.com/"),
	"www: archlinux":                  OpenWeb("https://wiki.archlinux.org"),
	"www: arxiv":                      OpenWeb("https://arxiv.org"),
	"www: bilibili":                   OpenWeb("https://www.bilibili.com"),
	"www: bing":                       OpenWeb("https://cn.bing.com"),
	"www: github":                     OpenWeb("https://github.com"),
	"www: github repos":               OpenWeb("https://github.com/zetatez?tab=repositories"),
	"www: gitee":                      OpenWeb("https://gitee.com"),
	"www: gitee repos":                OpenWeb("https://gitee.com/zetatez/projects"),
	"www: instagram":                  OpenWeb("https://www.instagram.com/explore/"),
	"www: jd":                         OpenWeb("https://www.jd.com"),
	"www: mirror aliyun":              OpenWeb("https://developer.aliyun.com/mirror"),
	"www: news":                       OpenWeb("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"www: scholar":                    OpenWeb("https://scholar.google.com"),
	"www: suckless":                   OpenWeb("https://dwm.suckless.org"),
	"www: tv cctv5":                   OpenWeb("https://tv.cctv.com/live/cctv5"),
	"www: twitter":                    OpenWeb("https://twitter.com/home"),
	"www: wolframalpha":               OpenWeb("https://www.wolframalpha.com"),
	"www: youtube":                    OpenWeb("https://www.youtube.com"),
	"dev: regex":                      OpenWeb("https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference"),
	"dev: runoob":                     OpenWeb("https://www.runoob.com"),
	"dev: css":                        OpenWeb("https://www.runoob.com/css3/css3-tutorial.html"),
	"dev: design pattern":             OpenWeb("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"dev: docker":                     OpenWeb("https://www.runoob.com/docker/docker-tutorial.html"),
	"dev: html":                       OpenWeb("https://www.runoob.com/html/html5-intro.html"),
	"dev: javascript":                 OpenWeb("https://www.runoob.com/js/js-tutorial.html"),
	"dev: maven":                      OpenWeb("https://www.runoob.com/maven/maven-tutorial.html"),
	"dev: mongo":                      OpenWeb("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"dev: nodejs":                     OpenWeb("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"dev: react":                      OpenWeb("https://www.runoob.com/react/react-tutorial.html"),
	"dev: redis":                      OpenWeb("https://www.runoob.com/redis/redis-tutorial.html"),
	"dev: typescript":                 OpenWeb("https://www.runoob.com/typescript/ts-tutorial.html"),
	"dev: vue":                        OpenWeb("https://www.runoob.com/vue3/vue3-tutorial.html"),
	"wiki: consul":                    OpenWeb("https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"wiki: scala sbt":                 OpenWeb("https://www.scala-sbt.org"),
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
	sort.Strings(list)
	content, err := sugar.Choose("search: ", list)
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
		OpenWeb(content)()
		return
	default:
		SearchFromWeb(content)
	}
}
