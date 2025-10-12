package plugins

import (
	"sort"

	"cmds/utils"
)

var ActionMap = map[string]func(){
	"format json":                      FormatJson,
	"format sql":                       FormatSql,
	"format yaml":                      FormatYaml,
	"get current datetime":             GetCurrentDatetime,
	"get current unix sec":             GetCurrentUnixSec,
	"get host name":                    GetHostName,
	"get ip address":                   GetIPAddress,
	"handle copied":                    HandleCopied,
	"jump to code from log":            JumpToCodeFromLog,
	"note script":                      NoteScripts,
	"note todo":                        NoteToDo,
	"note diary":                       NoteDiary,
	"note flash card":                  NoteFlashCard,
	"note timeline":                    NoteTimeline,
	"ssh to":                           SshTo,
	"search books online":              SearchBooksOnline,
	"search videos online":             SearchVideosOnline,
	"sys bluetooth connect":            SysBlueToothConnect,
	"sys bluetooth disconnect":         SysBlueToothDisconnect,
	"sys bluetooth scan and connect":   SysBlueToothScanAndConnect,
	"sys screen":                       SysScreen,
	"sys shortcuts":                    SysShortcuts,
	"sys toggle keyboard light":        SysToggleKeyboardLight,
	"sys wifi connect":                 SysWifiConnect,
	"toggle addressbook":               ToggleAddressbook,
	"toggle calendar scheduling today": ToggleCalendarSchedulingToday,
	"toggle calendar scheduling":       ToggleCalendarScheduling,
	"toggle calendar":                  ToggleCalendar,
	"toggle clipmenu":                  ToggleClipmenu,
	"toggle flameshot":                 ToggleFlameshot,
	"toggle inkscape":                  ToggleInkscape,
	"toggle irssi":                     ToggleIrssi,
	"toggle julia":                     ToggleJulia,
	"toggle krita":                     ToggleKrita,
	"toggle lazydocker":                ToggleLazyDocker,
	"toggle lazygit":                   ToggleLazyGit,
	"toggle lua":                       ToggleLua,
	"toggle music":                     ToggleMusic,
	"toggle music-net-cloud":           ToggleMusicNetCloud,
	"toggle mutt":                      ToggleMutt,
	"toggle obsidian":                  ToggleObsidian,
	"toggle passmenu":                  TogglePassmenu,
	"toggle python":                    TogglePython,
	"toggle rec-audio":                 ToggleRecAudio,
	"toggle rec-screen":                ToggleRecScreen,
	"toggle rec-webcam":                ToggleRecWebcam,
	"toggle redshift":                  ToggleRedShift,
	"toggle screenkey":                 ToggleScreenKey,
	"toggle show":                      ToggleShow,
	"toggle sublime":                   ToggleSublime,
	"toggle top":                       ToggleTop,
	"toggle tty_clock":                 ToggleTTYClock,
	"toggle xournal":                   ToggleXournal,
	"toggle yazi":                      ToggleYazi,
	"transform datetime to unix sec":   TransformDatetime2UnixSec,
	"transform unix sec to datetime":   TransformUnixSec2DateTime,
	"web ai chatgpt":                   OpenUrlWithQutebrowser("https://chatgpt.com/"),
	"web ai doubao":                    OpenUrlWithQutebrowser("https://www.doubao.com/"),
	"web email 163":                    OpenUrlWithQutebrowser("https://mail.163.com"),
	"web email gmail":                  OpenUrlWithQutebrowser("https://accounts.google.com/b/0/AddMailService"),
	"web email outlook":                OpenUrlWithQutebrowser("https://outlook.live.com/mail"),
	"web github":                       OpenUrlWithQutebrowser("https://github.com"),
	"web github repos":                 OpenUrlWithQutebrowser("https://github.com/zetatez?tab=repositories"),
	"web github gist":                  OpenUrlWithQutebrowser("https://gist.github.com/"),
	"web fin antfin":                   OpenUrlWithQutebrowser("https://caifu.antfin.com/"),
	"web fin eastmoney fund":           OpenUrlWithQutebrowser("https://fund.eastmoney.com/"),
	"web fin eastmoney":                OpenUrlWithQutebrowser("https://www.eastmoney.com/"),
	"web fin futu":                     OpenUrlWithQutebrowser("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"web fin tonghuashun":              OpenUrlWithQutebrowser("https://www.10jqka.com.cn/"),
	"web fin xueqiu":                   OpenUrlWithQutebrowser("https://xueqiu.com/"),
	"web google translate auto to en":  OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=en"),
	"web google translate auto to zh":  OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web life 12306":                   OpenUrlWithQutebrowser("https://www.12306.cn/index/"),
	"web life ctrip":                   OpenUrlWithQutebrowser("https://www.ctrip.com/"),
	"web life da.zhong.dian.ping":      OpenUrlWithQutebrowser("https://www.dianping.com/"),
	"web life jd":                      OpenUrlWithQutebrowser("https://www.jd.com"),
	"web life gaode map":               OpenUrlWithQutebrowser("https://ditu.amap.com/"),
	"web life google map":              OpenUrlWithQutebrowser("https://www.google.com/maps/place/shanghai"),
	"web life meituan":                 OpenUrlWithQutebrowser("https://www.meituan.com/"),
	"web mirror aliyun":                OpenUrlWithQutebrowser("https://developer.aliyun.com/mirror"),
	"web office feishu docs":           OpenUrlWithQutebrowser("https://docs.feishu.cn/docs/"),
	"web office feishu meeting":        OpenUrlWithQutebrowser("https://meeting.feishu.cn/"),
	"web ref archlinux":                OpenUrlWithQutebrowser("https://wiki.archlinux.org"),
	"web ref consul":                   OpenUrlWithQutebrowser("https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"web ref data-structures":          OpenUrlWithQutebrowser("https://www.runoob.com/data-structures/data-structures-tutorial.html"),
	"web ref db mongodb":               OpenUrlWithQutebrowser("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"web ref db mysql":                 OpenUrlWithQutebrowser("https://www.runoob.com/mysql/mysql-tutorial.html"),
	"web ref db postgresql":            OpenUrlWithQutebrowser("https://www.runoob.com/postgresql/postgresql-tutorial.html"),
	"web ref db redis":                 OpenUrlWithQutebrowser("https://www.runoob.com/redis/redis-tutorial.html"),
	"web ref db sqlite":                OpenUrlWithQutebrowser("https://www.runoob.com/sqlite/sqlite-tutorial.html"),
	"web ref design pattern":           OpenUrlWithQutebrowser("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"web ref docker":                   OpenUrlWithQutebrowser("https://www.runoob.com/docker/docker-tutorial.html"),
	"web ref git":                      OpenUrlWithQutebrowser("https://www.runoob.com/git/git-tutorial.html"),
	"web ref lang css":                 OpenUrlWithQutebrowser("https://www.runoob.com/css3/css3-tutorial.html"),
	"web ref lang flutter":             OpenUrlWithQutebrowser("https://flutter.dev/docs"),
	"web ref lang go":                  OpenUrlWithQutebrowser("https://www.runoob.com/go/go-tutorial.html"),
	"web ref lang html":                OpenUrlWithQutebrowser("https://www.runoob.com/html/html5-intro.html"),
	"web ref lang java":                OpenUrlWithQutebrowser("https://docs.oracle.com/en/java/"),
	"web ref lang javascript":          OpenUrlWithQutebrowser("https://www.runoob.com/js/js-tutorial.html"),
	"web ref lang julia":               OpenUrlWithQutebrowser("https://www.runoob.com/julia/julia-tutorial.html"),
	"web ref lang kotlin":              OpenUrlWithQutebrowser("https://kotlinlang.org/docs/home.html"),
	"web ref lang maven":               OpenUrlWithQutebrowser("https://www.runoob.com/maven/maven-tutorial.html"),
	"web ref lang nodejs":              OpenUrlWithQutebrowser("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"web ref lang python":              OpenUrlWithQutebrowser("https://docs.python.org/zh-cn/"),
	"web ref lang rust":                OpenUrlWithQutebrowser("https://www.runoob.com/rust/rust-tutorial.html"),
	"web ref lang scala":               OpenUrlWithQutebrowser("https://www.runoob.com/scala/scala-tutorial.html"),
	"web ref lang scala-sbt":           OpenUrlWithQutebrowser("https://www.scala-sbt.org"),
	"web ref lang sql":                 OpenUrlWithQutebrowser("https://www.runoob.com/sql/sql-tutorial.html"),
	"web ref lang swift":               OpenUrlWithQutebrowser("https://developer.apple.com/xcode/swiftui/"),
	"web ref lang typescript":          OpenUrlWithQutebrowser("https://www.runoob.com/typescript/ts-tutorial.html"),
	"web ref lang vue":                 OpenUrlWithQutebrowser("https://www.runoob.com/vue3/vue3-tutorial.html"),
	"web ref regex":                    OpenUrlWithQutebrowser("https://www.runoob.com/regexp/regexp-tutorial.html"),
	"web ref runoob":                   OpenUrlWithQutebrowser("https://www.runoob.com"),
	"web study leetcode":               OpenUrlWithQutebrowser("https://leetcode.cn/problemset/"),
	"web study wolframalpha":           OpenUrlWithQutebrowser("https://www.wolframalpha.com"),
	"web study arxiv":                  OpenUrlWithQutebrowser("https://arxiv.org"),
	"web study geekbang":               OpenUrlWithQutebrowser("https://time.geekbang.org/"),
	"web study scholar":                OpenUrlWithQutebrowser("https://scholar.google.com"),
	"web study shanbay":                OpenUrlWithQutebrowser("https://www.shanbay.com/"),
	"web study youdao":                 OpenUrlWithQutebrowser("https://dict.youdao.com/"),
	"web tool online photoshop":        OpenUrlWithQutebrowser("https://www.photopea.com/"),
	"web tool wechat file help":        OpenUrlWithQutebrowser("https://filehelper.weixin.qq.com/"),
	"web social douban":                OpenUrlWithQutebrowser("https://www.douban.com/"),
	"web social instagram":             OpenUrlWithQutebrowser("https://www.instagram.com/explore/"),
	"web social wechat":                OpenUrlWithQutebrowser("https://web.wechat.com/"),
	"web video bilibili":               OpenUrlWithQutebrowser("https://www.bilibili.com"),
	"web video cctv5":                  OpenUrlWithQutebrowser("https://tv.cctv.com/live/cctv5"),
	"web video iqiyi":                  OpenUrlWithQutebrowser("https://www.iqiyi.com/"),
	"web video youtube":                OpenUrlWithQutebrowser("https://www.youtube.com"),
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
	content, err := utils.Choose("search: ", list)
	if err != nil {
		return
	}
	f, ok := ActionMap[content]
	switch {
	case ok:
		f()
		return
	case utils.Exists(content) && utils.IsFile(content):
		utils.Lazy("open", content)
		return
	case utils.IsURL(content):
		OpenUrlWithQutebrowser(content)()
		return
	default:
		SearchFromWeb(content)
	}
}
