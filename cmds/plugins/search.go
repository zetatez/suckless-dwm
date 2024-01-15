package plugins

import (
	"sort"

	"cmds/sugar"
)

var ActionMap = map[string]func(){
	"launchapp: baidunetdisk":        LaunchApp("baidunetdisk"),
	"launchapp: chrome":              LaunchApp("chrome --proxy-server=socks5://127.0.0.1:7891"),
	"launchapp: dingtalk":            LaunchApp("dingtalk"),
	"launchapp: edge":                LaunchApp("microsoft-edge-stable"),
	"launchapp: inkscape":            LaunchApp("inkscape"),
	"launchapp: krita":               LaunchApp("krita"),
	"launchapp: netease-cloud-music": LaunchApp("netease-cloud-music"),
	"launchapp: obsidian":            LaunchApp("obsidian"),
	"launchapp: passmenu":            LaunchApp("passmenu"),
	"launchapp: passmmenu":           LaunchApp("passmmenu"),
	"launchapp: scribus":             LaunchApp("scribus"),
	"launchapp: slack":               LaunchApp("slack"),
	"launchapp: subl":                LaunchApp("subl"),
	"launchapp: wechat-uos":          LaunchApp("wechat-uos"),
	"launchapp: wemeet":              LaunchApp("wemeet"),
	"launchapp: wps":                 LaunchApp("wps"),
	"launchapp: xournalpp":           LaunchApp("xournalpp"),
	"launchapp: zoom":                LaunchApp("zoom"),
	"toggle: address-book":           ToggleAddressbook,
	"toggle: bluetooth":              ToggleBlueTooth,
	"toggle: calendar-day":           ToggleCalendarDay,
	"toggle: calendar-week":          ToggleCalendarWeek,
	"toggle: chrome":                 ToggleChrome,
	"toggle: clipmenu":               ToggleClipmenu,
	"toggle: diary":                  ToggleDiary,
	"toggle: flameshot":              ToggleFlameshot,
	"toggle: inkscape":               ToggleInkscape,
	"toggle: irssi":                  ToggleIrssi,
	"toggle: joshuto":                ToggleJoshuto,
	"toggle: keyboard-light":         ToggleKeyboardLight,
	"toggle: krita":                  ToggleKrita,
	"toggle: lazydocker":             ToggleLazyDocker,
	"toggle: music":                  ToggleMusic,
	"toggle: music-net-cloud":        ToggleMusicNetCloud,
	"toggle: mutt":                   ToggleMutt,
	"toggle: passmenu":               TogglePassmenu,
	"toggle: julia":                  ToggleJulia,
	"toggle: python":                 TogglePython,
	"toggle: scala":                  ToggleScala,
	"toggle: lua":                    ToggleLua,
	"toggle: recaudio":               ToggleRecAudio,
	"toggle: recvideo":               ToggleRecVideo,
	"toggle: redshift":               ToggleRedShift,
	"toggle: screen":                 ToggleScreen,
	"toggle: screenkey":              ToggleScreenKey,
	"toggle: show":                   ToggleShow,
	"toggle: sublime":                ToggleSublime,
	"toggle: sys-short-cuts":         ToggleSysShortcuts,
	"toggle: top":                    ToggleTop,
	"toggle: wallpaper":              ToggleWallpaper,
	"toggle: wechat":                 ToggleWechat,
	"toggle: xournal":                ToggleXournal,
	"workflow: ssh":                  SSH,
	"workflow: format-sql":           FormatSql,
	"workflow: format-json":          FormatSql,
	"workflow: current-datetime":     CurrentDatetime,
	"workflow: current-unix-sec":     CurrentUnixSec,
	"workflow: handle-copied":        HandleCopied,
	"workflow: umount xyz":           UmountXYZ,
	"workflow: wifi-connect":         WifiConnect,
	"website: arch wiki":             Website("https://wiki.archlinux.org"),
	"website: arxiv":                 Website("https://arxiv.org"),
	"website: bili":                  Website("https://www.bilibili.com"),
	"website: bing":                  Website("https://cn.bing.com"),
	"website: cctv5":                 Website("https://tv.cctv.com/live/cctv5"),
	"website: github":                Website("https://github.com/zetatez?tab=repositories"),
	"website: mall":                  Website("https://www.jd.com"),
	"website: map":                   Website("https://www.google.com/maps/place/shanghai"),
	"website: mirror":                Website("https://developer.aliyun.com/mirror"),
	"website: news":                  Website("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"website: ocr":                   Website("http://ocr.space"),
	"website: regex":                 Website("https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference"),
	"website: scholar":               Website("https://scholar.google.com"),
	"website: suckless":              Website("https://dwm.suckless.org"),
	"website: translation":           Website("https://translate.google.com/?hl=zh-CN"),
	"website: wolframalpha":          Website("https://www.wolframalpha.com"),
	"website: youtube":               Website("https://www.youtube.com"),
	"website: runoob":                Website("https://www.runoob.com"),
	"website: css":                   Website("https://www.runoob.com/css3/css3-tutorial.html"),
	"website: design pattern":        Website("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"website: docker":                Website("https://www.runoob.com/docker/docker-tutorial.html"),
	"website: html":                  Website("https://www.runoob.com/html/html5-intro.html"),
	"website: javascript":            Website("https://www.runoob.com/js/js-tutorial.html"),
	"website: maven":                 Website("https://www.runoob.com/maven/maven-tutorial.html"),
	"website: mongo":                 Website("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"website: nodejs":                Website("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"website: react":                 Website("https://www.runoob.com/react/react-tutorial.html"),
	"website: redis":                 Website("https://www.runoob.com/redis/redis-tutorial.html"),
	"website: typescript":            Website("https://www.runoob.com/typescript/ts-tutorial.html"),
	"website: vue":                   Website("https://www.runoob.com/vue3/vue3-tutorial.html"),
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
		Website(content)()
		return
	default:
		SearchFromWeb(content)
	}
}
