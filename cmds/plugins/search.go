package plugins

import (
	"sort"

	"cmds/sugar"
)

var ActionMap = map[string]func(){
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
	"transform date time to unix sec": TransformDatetime2UnixSec,
	"transform unix sec to date time": TransformUnixSec2DateTime,
	"note diary":                      NoteDiary,
	"note timeline":                   NoteTimeline,
	"note flash card":                 NoteFlashCard,
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
	"toggle xournal":                  ToggleXournal,
	"jump to code from log":           JumpToCodeFromLog,
	"web google map":                  OpenWeb("--proxy-server="+ProxyServer, "https://www.google.com/maps/place/shanghai"),
	"web google translate auto to en": OpenWeb("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=en"),
	"web google translate auto to zh": OpenWeb("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web email gmail":                 OpenWeb("--proxy-server="+ProxyServer, "https://accounts.google.com/b/0/AddMailService"),
	"web email outlook":               OpenWeb("--proxy-server="+ProxyServer, "https://outlook.live.com/mail"),
	"web email 163":                   OpenWeb("", "https://mail.163.com"),
	"web gitee repos":                 OpenWeb("", "https://gitee.com/zetatez/projects"),
	"web gitee":                       OpenWeb("", "https://gitee.com"),
	"web github gist":                 OpenWeb("--proxy-server="+ProxyServer, "https://gist.github.com/"),
	"web github repos":                OpenWeb("--proxy-server="+ProxyServer, "https://github.com/zetatez?tab=repositories"),
	"web github":                      OpenWeb("--proxy-server="+ProxyServer, "https://github.com"),
	"web learning leetcode":           OpenWeb("", "https://leetcode.cn/problemset/"),
	"web learning wolframalpha":       OpenWeb("", "https://www.wolframalpha.com"),
	"web mall jd":                     OpenWeb("", "https://www.jd.com"),
	"web paper arxiv":                 OpenWeb("--proxy-server="+ProxyServer, "https://arxiv.org"),
	"web paper scholar":               OpenWeb("--proxy-server="+ProxyServer, "https://scholar.google.com"),
	"web search engine bing":          OpenWeb("", "https://cn.bing.com"),
	"web social instagram":            OpenWeb("--proxy-server="+ProxyServer, "https://www.instagram.com/explore/"),
	"web social twitter":              OpenWeb("--proxy-server="+ProxyServer, "https://twitter.com/home"),
	"web social wechat":               OpenWeb("", "https://web.wechat.com/"),
	"web tool chatgpt":                OpenWeb("--proxy-server="+ProxyServer, "https://chatgpt.com/"),
	"web tool mirror aliyun":          OpenWeb("", "https://developer.aliyun.com/mirror"),
	"web tool news":                   OpenWeb("", "https://news.futunn.com/en/main/live?lang=zh-CN"),
	"web tool wechat file help":       OpenWeb("", "https://filehelper.weixin.qq.com/"),
	"web video bilibili":              OpenWeb("", "https://www.bilibili.com"),
	"web video cctv5":                 OpenWeb("", "https://tv.cctv.com/live/cctv5"),
	"web video youtube":               OpenWeb("--proxy-server="+ProxyServer, "https://www.youtube.com"),
	"web reference archlinux":         OpenWeb("--proxy-server="+ProxyServer, "https://wiki.archlinux.org"),
	"web reference consul":            OpenWeb("--proxy-server="+ProxyServer, "https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"web reference data-structures":   OpenWeb("", "https://www.runoob.com/data-structures/data-structures-tutorial.html"),
	"web reference db mongodb":        OpenWeb("", "https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"web reference db mysql":          OpenWeb("", "https://www.runoob.com/mysql/mysql-tutorial.html"),
	"web reference db postgresql":     OpenWeb("", "https://www.runoob.com/postgresql/postgresql-tutorial.html"),
	"web reference db redis":          OpenWeb("", "https://www.runoob.com/redis/redis-tutorial.html"),
	"web reference db sqlite":         OpenWeb("", "https://www.runoob.com/sqlite/sqlite-tutorial.html"),
	"web reference design pattern":    OpenWeb("", "https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"web reference docker":            OpenWeb("", "https://www.runoob.com/docker/docker-tutorial.html"),
	"web reference git":               OpenWeb("", "https://www.runoob.com/git/git-tutorial.html"),
	"web reference lang css":          OpenWeb("", "https://www.runoob.com/css3/css3-tutorial.html"),
	"web reference lang go":           OpenWeb("", "https://www.runoob.com/go/go-tutorial.html"),
	"web reference lang html":         OpenWeb("", "https://www.runoob.com/html/html5-intro.html"),
	"web reference lang javascript":   OpenWeb("", "https://www.runoob.com/js/js-tutorial.html"),
	"web reference lang julia":        OpenWeb("", "https://www.runoob.com/julia/julia-tutorial.html"),
	"web reference lang maven":        OpenWeb("", "https://www.runoob.com/maven/maven-tutorial.html"),
	"web reference lang nodejs":       OpenWeb("", "https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"web reference lang python":       OpenWeb("", "https://www.runoob.com/python3/python3-tutorial.html"),
	"web reference lang rust":         OpenWeb("", "https://www.runoob.com/rust/rust-tutorial.html"),
	"web reference lang scala":        OpenWeb("", "https://www.runoob.com/scala/scala-tutorial.html"),
	"web reference lang scala-sbt":    OpenWeb("", "https://www.scala-sbt.org"),
	"web reference lang sql":          OpenWeb("", "https://www.runoob.com/sql/sql-tutorial.html"),
	"web reference lang typescript":   OpenWeb("", "https://www.runoob.com/typescript/ts-tutorial.html"),
	"web reference lang vue":          OpenWeb("", "https://www.runoob.com/vue3/vue3-tutorial.html"),
	"web reference regex":             OpenWeb("", "https://www.runoob.com/regexp/regexp-tutorial.html"),
	"web reference runoob":            OpenWeb("", "https://www.runoob.com"),
	"web reference suckless":          OpenWeb("--proxy-server="+ProxyServer, "https://dwm.suckless.org"),
	"web sound 雨水-01":                 OpenWeb("--proxy-server="+ProxyServer, "https://www.youtube.com/watch?v=O8o3T01reS0&ab_channel=%E8%87%AA%E7%84%B6%E9%9F%B3%E6%A8%82"),
	"web sound 溪流-01":                 OpenWeb("--proxy-server="+ProxyServer, "https://www.youtube.com/watch?v=YjUkT-Ufrv8&ab_channel=%E7%99%92%E3%81%97%E3%81%AE%E6%B0%B4%E3%81%AE%E9%9F%B3ch"),
	"web sound 溪流-02":                 OpenWeb("", "https://www.bilibili.com/video/BV1TG4y1a7U1/?spm_id_from=333.999.0.0&vd_source=869c7edafe114294eee747bd802bd2dd"),
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
		OpenWeb("--proxy-server="+ProxyServer, content)()
		return
	default:
		SearchFromWeb(content)
	}
}
