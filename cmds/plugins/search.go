package plugins

import (
	"fmt"
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
	"google map":                      OpenWeb("https://www.google.com/maps/place/shanghai"),
	"google translate auto to en":     OpenWeb("https://translate.google.com/?sl=auto&tl=en"),
	"google translate auto to zh":     OpenWeb("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"email 163":                       OpenWeb("https://mail.163.com"),
	"email gmail":                     OpenWeb("https://accounts.google.com/b/0/AddMailService"),
	"email outlook":                   OpenWeb("https://outlook.live.com/mail"),
	"transform date time to unix sec": TransformDatetime2UnixSec,
	"transform unix sec to date time": TransformUnixSec2DateTime,
	"note diary":                      NoteDiary,
	"note timeline":                   NoteTimeline,
	"note flash card":                 NoteFlashCard,
	"launch baidudisknet":             LaunchApp("baidudisknet"),
	"launch chrome":                   LaunchApp(fmt.Sprintf("chrome --proxy-server=%s", ProxyServer)),
	"launch edge":                     LaunchApp(fmt.Sprintf("edge --proxy-server=%s", ProxyServer)),
	"launch dingtalk":                 LaunchApp("dingtalk"),
	"launch inkscape":                 LaunchApp("inkscape"),
	"launch krita":                    LaunchApp("krita"),
	"launch netease-cloud-music":      LaunchApp("netease-cloud-music"),
	"launch obsidian":                 LaunchApp("obsidian"),
	"launch passmenu":                 LaunchApp("passmenu"),
	"launch passmmenu":                LaunchApp("passmmenu"),
	"launch scribus":                  LaunchApp("scribus"),
	"launch slack":                    LaunchApp("slack"),
	"launch subl":                     LaunchApp("subl"),
	"launch wechat":                   LaunchApp("wechat"),
	"launch wemeet":                   LaunchApp("wemeet"),
	"launch wps":                      LaunchApp("wps"),
	"launch xournalpp":                LaunchApp("xournalpp"),
	"launch zoom":                     LaunchApp("zoom"),
	"toggle address-book":             ToggleAddressbook,
	"toggle bluetooth":                ToggleBlueTooth,
	"toggle calendar today schedule":  ToggleCalendarTodaySchedule,
	"toggle calendar scheduling":      ToggleCalendarScheduling,
	"toggle clipmenu":                 ToggleClipmenu,
	"toggle flameshot":                ToggleFlameshot,
	"toggle inkscape":                 ToggleInkscape,
	"toggle irssi":                    ToggleIrssi,
	"toggle joshuto":                  ToggleJoshuto,
	"toggle julia":                    ToggleJulia,
	"toggle keyboard-light":           ToggleKeyboardLight,
	"toggle krita":                    ToggleKrita,
	"toggle lazydocker":               ToggleLazyDocker,
	"toggle lua":                      ToggleLua,
	"toggle music":                    ToggleMusic,
	"toggle music-net-cloud":          ToggleMusicNetCloud,
	"toggle mutt":                     ToggleMutt,
	"toggle passmenu":                 TogglePassmenu,
	"toggle python":                   TogglePython,
	"toggle rec-audio":                ToggleRecAudio,
	"toggle rec-screen":               ToggleRecScreen,
	"toggle rec-webcam":               ToggleRecWebcam,
	"toggle redshift":                 ToggleRedShift,
	"toggle scala":                    ToggleScala,
	"toggle screen":                   ToggleScreen,
	"toggle screenkey":                ToggleScreenKey,
	"toggle show":                     ToggleShow,
	"toggle sublime":                  ToggleSublime,
	"toggle sys-short-cuts":           ToggleSysShortcuts,
	"toggle top":                      ToggleTop,
	"toggle wallpaper":                ToggleWallpaper,
	"toggle wechat":                   ToggleWechat,
	"toggle xournal":                  ToggleXournal,
	"jump to code from log":           JumpToCodeFromLog,
	"web gitee repos":                 OpenWeb("https://gitee.com/zetatez/projects"),
	"web gitee":                       OpenWeb("https://gitee.com"),
	"web github gist":                 OpenWeb("https://gist.github.com/"),
	"web github repos":                OpenWeb("https://github.com/zetatez?tab=repositories"),
	"web github":                      OpenWeb("https://github.com"),
	"web learning leetcode":           OpenWeb("https://leetcode.cn/problemset/"),
	"web learning wolframalpha":       OpenWeb("https://www.wolframalpha.com"),
	"web mall jd":                     OpenWeb("https://www.jd.com"),
	"web paper arxiv":                 OpenWeb("https://arxiv.org"),
	"web paper scholar":               OpenWeb("https://scholar.google.com"),
	"web search engine bing":          OpenWeb("https://cn.bing.com"),
	"web social instagram":            OpenWeb("https://www.instagram.com/explore/"),
	"web social twitter":              OpenWeb("https://twitter.com/home"),
	"web social wechat":               OpenWeb("https://web.wechat.com/"),
	"web tool chatgpt":                OpenWeb("https://chatgpt.com/"),
	"web tool mirror aliyun":          OpenWeb("https://developer.aliyun.com/mirror"),
	"web tool news":                   OpenWeb("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"web tool wechat file help":       OpenWeb("https://filehelper.weixin.qq.com/"),
	"web video bilibili":              OpenWeb("https://www.bilibili.com"),
	"web video cctv5":                 OpenWeb("https://tv.cctv.com/live/cctv5"),
	"web video youtube":               OpenWeb("https://www.youtube.com"),
	"web reference archlinux":         OpenWeb("https://wiki.archlinux.org"),
	"web reference consul":            OpenWeb("https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"web reference css":               OpenWeb("https://www.runoob.com/css3/css3-tutorial.html"),
	"web reference design pattern":    OpenWeb("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"web reference docker":            OpenWeb("https://www.runoob.com/docker/docker-tutorial.html"),
	"web reference html":              OpenWeb("https://www.runoob.com/html/html5-intro.html"),
	"web reference javascript":        OpenWeb("https://www.runoob.com/js/js-tutorial.html"),
	"web reference maven":             OpenWeb("https://www.runoob.com/maven/maven-tutorial.html"),
	"web reference mongo":             OpenWeb("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"web reference nodejs":            OpenWeb("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"web reference react":             OpenWeb("https://www.runoob.com/react/react-tutorial.html"),
	"web reference redis":             OpenWeb("https://www.runoob.com/redis/redis-tutorial.html"),
	"web reference regex":             OpenWeb("https://learn.microsoft.com/zh-cn/dotnet/standard/base-types/regular-expression-language-quick-reference"),
	"web reference runoob":            OpenWeb("https://www.runoob.com"),
	"web reference scala sbt":         OpenWeb("https://www.scala-sbt.org"),
	"web reference suckless":          OpenWeb("https://dwm.suckless.org"),
	"web reference typescript":        OpenWeb("https://www.runoob.com/typescript/ts-tutorial.html"),
	"web reference vue":               OpenWeb("https://www.runoob.com/vue3/vue3-tutorial.html"),
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
