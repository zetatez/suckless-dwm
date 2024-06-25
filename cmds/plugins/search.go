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
	"google map":                      OpenWeb("https://www.google.com/maps/place/shanghai"),
	"google translate auto to en":     OpenWeb("https://translate.google.com/?sl=auto&tl=en"),
	"google translate auto to zh":     OpenWeb("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"email gmail":                     OpenWeb("https://accounts.google.com/b/0/AddMailService"),
	"email outlook":                   OpenWeb("https://outlook.live.com/mail"),
	"email 163":                       OpenWeb("https://mail.163.com"),
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
	"web reference data-structures":   OpenWeb("https://www.runoob.com/data-structures/data-structures-tutorial.html"),
	"web reference db mongodb":        OpenWeb("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"web reference db mysql":          OpenWeb("https://www.runoob.com/mysql/mysql-tutorial.html"),
	"web reference db postgresql":     OpenWeb("https://www.runoob.com/postgresql/postgresql-tutorial.html"),
	"web reference db redis":          OpenWeb("https://www.runoob.com/redis/redis-tutorial.html"),
	"web reference db sqlite":         OpenWeb("https://www.runoob.com/sqlite/sqlite-tutorial.html"),
	"web reference design pattern":    OpenWeb("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"web reference docker":            OpenWeb("https://www.runoob.com/docker/docker-tutorial.html"),
	"web reference git":               OpenWeb("https://www.runoob.com/git/git-tutorial.html"),
	"web reference lang css":          OpenWeb("https://www.runoob.com/css3/css3-tutorial.html"),
	"web reference lang go":           OpenWeb("https://www.runoob.com/go/go-tutorial.html"),
	"web reference lang html":         OpenWeb("https://www.runoob.com/html/html5-intro.html"),
	"web reference lang javascript":   OpenWeb("https://www.runoob.com/js/js-tutorial.html"),
	"web reference lang julia":        OpenWeb("https://www.runoob.com/julia/julia-tutorial.html"),
	"web reference lang maven":        OpenWeb("https://www.runoob.com/maven/maven-tutorial.html"),
	"web reference lang nodejs":       OpenWeb("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"web reference lang python":       OpenWeb("https://www.runoob.com/python3/python3-tutorial.html"),
	"web reference lang rust":         OpenWeb("https://www.runoob.com/rust/rust-tutorial.html"),
	"web reference lang scala":        OpenWeb("https://www.runoob.com/scala/scala-tutorial.html"),
	"web reference lang scala-sbt":    OpenWeb("https://www.scala-sbt.org"),
	"web reference lang sql":          OpenWeb("https://www.runoob.com/sql/sql-tutorial.html"),
	"web reference lang typescript":   OpenWeb("https://www.runoob.com/typescript/ts-tutorial.html"),
	"web reference lang vue":          OpenWeb("https://www.runoob.com/vue3/vue3-tutorial.html"),
	"web reference regex":             OpenWeb("https://www.runoob.com/regexp/regexp-tutorial.html"),
	"web reference runoob":            OpenWeb("https://www.runoob.com"),
	"web reference suckless":          OpenWeb("https://dwm.suckless.org"),
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
