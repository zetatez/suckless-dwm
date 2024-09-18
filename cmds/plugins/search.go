package plugins

import (
	"sort"

	"cmds/sugar"
)

var ActionMap = map[string]func(){
	"ssh-to":                           SshTo,
	"handle-copied":                    HandleCopied,
	"wifi-connect":                     WifiConnect,
	"format sql":                       FormatSql,
	"format json":                      FormatJson,
	"format yaml":                      FormatYaml,
	"search books online":              SearchBooksOnline,
	"search videos online":             SearchVideosOnline,
	"get current-datetime":             GetCurrentDatetime,
	"get current-unix-sec":             GetCurrentUnixSec,
	"get host name":                    GetHostName,
	"get ip":                           GetIP,
	"transform date time to unix sec":  TransformDatetime2UnixSec,
	"transform unix sec to date time":  TransformUnixSec2DateTime,
	"note diary":                       NoteDiary,
	"note timeline":                    NoteTimeline,
	"note flash card":                  NoteFlashCard,
	"toggle address-book":              ToggleAddressbook,
	"toggle bluetooth":                 ToggleBlueTooth,
	"toggle calendar scheduling today": ToggleCalendarSchedulingToday,
	"toggle calendar scheduling":       ToggleCalendarScheduling,
	"toggle clipmenu":                  ToggleClipmenu,
	"toggle flameshot":                 ToggleFlameshot,
	"toggle inkscape":                  ToggleInkscape,
	"toggle irssi":                     ToggleIrssi,
	"toggle joshuto":                   ToggleJoshuto,
	"toggle julia":                     ToggleJulia,
	"toggle keyboard-light":            ToggleKeyboardLight,
	"toggle krita":                     ToggleKrita,
	"toggle lazydocker":                ToggleLazyDocker,
	"toggle lua":                       ToggleLua,
	"toggle music":                     ToggleMusic,
	"toggle music-net-cloud":           ToggleMusicNetCloud,
	"toggle mutt":                      ToggleMutt,
	"toggle passmenu":                  TogglePassmenu,
	"toggle python":                    TogglePython,
	"toggle rec-audio":                 ToggleRecAudio,
	"toggle rec-screen":                ToggleRecScreen,
	"toggle rec-webcam":                ToggleRecWebcam,
	"toggle redshift":                  ToggleRedShift,
	"toggle scala":                     ToggleScala,
	"toggle screen":                    ToggleScreen,
	"toggle screenkey":                 ToggleScreenKey,
	"toggle show":                      ToggleShow,
	"toggle sublime":                   ToggleSublime,
	"toggle sys-short-cuts":            ToggleSysShortcuts,
	"toggle top":                       ToggleTop,
	"toggle wallpaper":                 ToggleWallpaper,
	"toggle xournal":                   ToggleXournal,
	"jump to code from log":            JumpToCodeFromLog,
	"url google map":                   ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.google.com/maps/place/shanghai"),
	"url google translate auto to en":  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=en"),
	"url google translate auto to zh":  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://translate.google.com/?sl=auto&tl=zh-CN"),
	"url email gmail":                  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://accounts.google.com/b/0/AddMailService"),
	"url email outlook":                ChromeOpenUrl("--proxy-server="+ProxyServer, "https://outlook.live.com/mail"),
	"url email 163":                    ChromeOpenUrl("", "https://mail.163.com"),
	"url gitee repos":                  ChromeOpenUrl("", "https://gitee.com/zetatez/projects"),
	"url gitee":                        ChromeOpenUrl("", "https://gitee.com"),
	"url github gist":                  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://gist.github.com/"),
	"url github repos":                 ChromeOpenUrl("--proxy-server="+ProxyServer, "https://github.com/zetatez?tab=repositories"),
	"url github":                       ChromeOpenUrl("--proxy-server="+ProxyServer, "https://github.com"),
	"url learning leetcode":            ChromeOpenUrl("", "https://leetcode.cn/problemset/"),
	"url learning wolframalpha":        ChromeOpenUrl("", "https://www.wolframalpha.com"),
	"url mall jd":                      ChromeOpenUrl("", "https://www.jd.com"),
	"url paper arxiv":                  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://arxiv.org"),
	"url paper scholar":                ChromeOpenUrl("--proxy-server="+ProxyServer, "https://scholar.google.com"),
	"url search engine bing":           ChromeOpenUrl("", "https://cn.bing.com"),
	"url social instagram":             ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.instagram.com/explore/"),
	"url social twitter":               ChromeOpenUrl("--proxy-server="+ProxyServer, "https://twitter.com/home"),
	"url social wechat":                ChromeOpenUrl("", "https://web.wechat.com/"),
	"url tool chatgpt":                 ChromeOpenUrl("--proxy-server="+ProxyServer, "https://chatgpt.com/"),
	"url tool mirror aliyun":           ChromeOpenUrl("", "https://developer.aliyun.com/mirror"),
	"url tool news":                    ChromeOpenUrl("", "https://news.futunn.com/en/main/live?lang=zh-CN"),
	"url tool wechat file help":        ChromeOpenUrl("", "https://filehelper.weixin.qq.com/"),
	"url video bilibili":               ChromeOpenUrl("", "https://www.bilibili.com"),
	"url video cctv5":                  ChromeOpenUrl("", "https://tv.cctv.com/live/cctv5"),
	"url video youtube":                ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.youtube.com"),
	"url reference archlinux":          ChromeOpenUrl("--proxy-server="+ProxyServer, "https://wiki.archlinux.org"),
	"url reference consul":             ChromeOpenUrl("--proxy-server="+ProxyServer, "https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"url reference data-structures":    ChromeOpenUrl("", "https://www.runoob.com/data-structures/data-structures-tutorial.html"),
	"url reference db mongodb":         ChromeOpenUrl("", "https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"url reference db mysql":           ChromeOpenUrl("", "https://www.runoob.com/mysql/mysql-tutorial.html"),
	"url reference db postgresql":      ChromeOpenUrl("", "https://www.runoob.com/postgresql/postgresql-tutorial.html"),
	"url reference db redis":           ChromeOpenUrl("", "https://www.runoob.com/redis/redis-tutorial.html"),
	"url reference db sqlite":          ChromeOpenUrl("", "https://www.runoob.com/sqlite/sqlite-tutorial.html"),
	"url reference design pattern":     ChromeOpenUrl("", "https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"url reference docker":             ChromeOpenUrl("", "https://www.runoob.com/docker/docker-tutorial.html"),
	"url reference git":                ChromeOpenUrl("", "https://www.runoob.com/git/git-tutorial.html"),
	"url reference lang css":           ChromeOpenUrl("", "https://www.runoob.com/css3/css3-tutorial.html"),
	"url reference lang go":            ChromeOpenUrl("", "https://www.runoob.com/go/go-tutorial.html"),
	"url reference lang html":          ChromeOpenUrl("", "https://www.runoob.com/html/html5-intro.html"),
	"url reference lang javascript":    ChromeOpenUrl("", "https://www.runoob.com/js/js-tutorial.html"),
	"url reference lang julia":         ChromeOpenUrl("", "https://www.runoob.com/julia/julia-tutorial.html"),
	"url reference lang maven":         ChromeOpenUrl("", "https://www.runoob.com/maven/maven-tutorial.html"),
	"url reference lang nodejs":        ChromeOpenUrl("", "https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"url reference lang python":        ChromeOpenUrl("", "https://www.runoob.com/python3/python3-tutorial.html"),
	"url reference lang rust":          ChromeOpenUrl("", "https://www.runoob.com/rust/rust-tutorial.html"),
	"url reference lang scala":         ChromeOpenUrl("", "https://www.runoob.com/scala/scala-tutorial.html"),
	"url reference lang scala-sbt":     ChromeOpenUrl("", "https://www.scala-sbt.org"),
	"url reference lang sql":           ChromeOpenUrl("", "https://www.runoob.com/sql/sql-tutorial.html"),
	"url reference lang typescript":    ChromeOpenUrl("", "https://www.runoob.com/typescript/ts-tutorial.html"),
	"url reference lang vue":           ChromeOpenUrl("", "https://www.runoob.com/vue3/vue3-tutorial.html"),
	"url reference regex":              ChromeOpenUrl("", "https://www.runoob.com/regexp/regexp-tutorial.html"),
	"url reference runoob":             ChromeOpenUrl("", "https://www.runoob.com"),
	"url reference suckless":           ChromeOpenUrl("--proxy-server="+ProxyServer, "https://dwm.suckless.org"),
	"url sound 雨水-01":                  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.youtube.com/watch?v=O8o3T01reS0&ab_channel=%E8%87%AA%E7%84%B6%E9%9F%B3%E6%A8%82"),
	"url sound 溪流-01":                  ChromeOpenUrl("--proxy-server="+ProxyServer, "https://www.youtube.com/watch?v=YjUkT-Ufrv8&ab_channel=%E7%99%92%E3%81%97%E3%81%AE%E6%B0%B4%E3%81%AE%E9%9F%B3ch"),
	"url sound 溪流-02":                  ChromeOpenUrl("", "https://www.bilibili.com/video/BV1TG4y1a7U1/?spm_id_from=333.999.0.0&vd_source=869c7edafe114294eee747bd802bd2dd"),
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
		ChromeOpenUrl("--proxy-server="+ProxyServer, content)()
		return
	default:
		SearchFromWeb(content)
	}
}
