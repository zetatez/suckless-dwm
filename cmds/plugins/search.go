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
	"web ai chatgpt":                   OpenUrlWithChrome("https://chatgpt.com/"),
	"web ai doubao":                    OpenUrlWithChrome("https://www.doubao.com/"),
	"web dev github Copilot":           OpenUrlWithChrome("https://github.com/features/copilot/"),
	"web dev stack Overflow":           OpenUrlWithChrome("https://stackoverflow.com/questions?tab=active"),
	"web email 163":                    OpenUrlWithChrome("https://mail.163.com"),
	"web email gmail":                  OpenUrlWithChrome("https://accounts.google.com/b/0/AddMailService"),
	"web email outlook":                OpenUrlWithChrome("https://outlook.live.com/mail"),
	"web finance antfin":               OpenUrlWithChrome("https://caifu.antfin.com/"),
	"web finance eastmoney fund":       OpenUrlWithChrome("https://fund.eastmoney.com/"),
	"web finance eastmoney":            OpenUrlWithChrome("https://www.eastmoney.com/"),
	"web finance futu":                 OpenUrlWithChrome("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"web finance tonghuashun":          OpenUrlWithChrome("https://www.10jqka.com.cn/"),
	"web finance xueqiu":               OpenUrlWithChrome("https://xueqiu.com/"),
	"web github gist":                  OpenUrlWithChrome("https://gist.github.com/"),
	"web github repos":                 OpenUrlWithChrome("https://github.com/zetatez?tab=repositories"),
	"web github":                       OpenUrlWithChrome("https://github.com"),
	"web google translate auto to en":  OpenUrlWithChrome("https://translate.google.com/?sl=auto&tl=en"),
	"web google translate auto to zh":  OpenUrlWithChrome("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web leetcode":                     OpenUrlWithChrome("https://leetcode.cn/problemset/"),
	"web life 12306":                   OpenUrlWithChrome("https://www.12306.cn/index/"),
	"web life ctrip":                   OpenUrlWithChrome("https://www.ctrip.com/"),
	"web life dianping":                OpenUrlWithChrome("https://www.dianping.com/"),
	"web life jd":                      OpenUrlWithChrome("https://www.jd.com"),
	"web life map gaode":               OpenUrlWithChrome("https://ditu.amap.com/"),
	"web life map google":              OpenUrlWithChrome("https://www.google.com/maps/place/shanghai"),
	"web life meituan":                 OpenUrlWithChrome("https://www.meituan.com/"),
	"web math wolframalpha":            OpenUrlWithChrome("https://www.wolframalpha.com"),
	"web mirror aliyun":                OpenUrlWithChrome("https://developer.aliyun.com/mirror"),
	"web office feishu docs":           OpenUrlWithChrome("https://docs.feishu.cn/docs/"),
	"web office feishu meeting":        OpenUrlWithChrome("https://meeting.feishu.cn/"),
	"web ref archlinux":                OpenUrlWithChrome("https://wiki.archlinux.org"),
	"web ref consul":                   OpenUrlWithChrome("https://developer.hashicorp.com/consul/docs?product_intent=consul"),
	"web ref data-structures":          OpenUrlWithChrome("https://www.runoob.com/data-structures/data-structures-tutorial.html"),
	"web ref db mongodb":               OpenUrlWithChrome("https://www.runoob.com/mongodb/mongodb-tutorial.html"),
	"web ref db mysql":                 OpenUrlWithChrome("https://www.runoob.com/mysql/mysql-tutorial.html"),
	"web ref db postgresql":            OpenUrlWithChrome("https://www.runoob.com/postgresql/postgresql-tutorial.html"),
	"web ref db redis":                 OpenUrlWithChrome("https://www.runoob.com/redis/redis-tutorial.html"),
	"web ref db sqlite":                OpenUrlWithChrome("https://www.runoob.com/sqlite/sqlite-tutorial.html"),
	"web ref design pattern":           OpenUrlWithChrome("https://www.runoob.com/design-pattern/design-pattern-tutorial.html"),
	"web ref docker":                   OpenUrlWithChrome("https://www.runoob.com/docker/docker-tutorial.html"),
	"web ref git":                      OpenUrlWithChrome("https://www.runoob.com/git/git-tutorial.html"),
	"web ref lang css":                 OpenUrlWithChrome("https://www.runoob.com/css3/css3-tutorial.html"),
	"web ref lang flutter":             OpenUrlWithChrome("https://flutter.dev/docs"),
	"web ref lang go":                  OpenUrlWithChrome("https://www.runoob.com/go/go-tutorial.html"),
	"web ref lang html":                OpenUrlWithChrome("https://www.runoob.com/html/html5-intro.html"),
	"web ref lang java":                OpenUrlWithChrome("https://docs.oracle.com/en/java/"),
	"web ref lang javascript":          OpenUrlWithChrome("https://www.runoob.com/js/js-tutorial.html"),
	"web ref lang julia":               OpenUrlWithChrome("https://www.runoob.com/julia/julia-tutorial.html"),
	"web ref lang kotlin":              OpenUrlWithChrome("https://kotlinlang.org/docs/home.html"),
	"web ref lang maven":               OpenUrlWithChrome("https://www.runoob.com/maven/maven-tutorial.html"),
	"web ref lang nodejs":              OpenUrlWithChrome("https://www.runoob.com/nodejs/nodejs-tutorial.html"),
	"web ref lang python":              OpenUrlWithChrome("https://docs.python.org/zh-cn/"),
	"web ref lang rust":                OpenUrlWithChrome("https://www.runoob.com/rust/rust-tutorial.html"),
	"web ref lang scala":               OpenUrlWithChrome("https://www.runoob.com/scala/scala-tutorial.html"),
	"web ref lang scala-sbt":           OpenUrlWithChrome("https://www.scala-sbt.org"),
	"web ref lang sql":                 OpenUrlWithChrome("https://www.runoob.com/sql/sql-tutorial.html"),
	"web ref lang swift":               OpenUrlWithChrome("https://developer.apple.com/xcode/swiftui/"),
	"web ref lang typescript":          OpenUrlWithChrome("https://www.runoob.com/typescript/ts-tutorial.html"),
	"web ref lang vue":                 OpenUrlWithChrome("https://www.runoob.com/vue3/vue3-tutorial.html"),
	"web ref regex":                    OpenUrlWithChrome("https://www.runoob.com/regexp/regexp-tutorial.html"),
	"web ref runoob":                   OpenUrlWithChrome("https://www.runoob.com"),
	"web social douban":                OpenUrlWithChrome("https://www.douban.com/"),
	"web social instagram":             OpenUrlWithChrome("https://www.instagram.com/explore/"),
	"web social wechat":                OpenUrlWithChrome("https://web.wechat.com/"),
	"web study arxiv":                  OpenUrlWithChrome("https://arxiv.org"),
	"web study duolingo":               OpenUrlWithChrome("https://www.duolingo.com/"),
	"web study geekbang":               OpenUrlWithChrome("https://time.geekbang.org/"),
	"web study scholar":                OpenUrlWithChrome("https://scholar.google.com"),
	"web study shanbay":                OpenUrlWithChrome("https://www.shanbay.com/"),
	"web study youdao":                 OpenUrlWithChrome("https://dict.youdao.com/"),
	"web tool online photoshop":        OpenUrlWithChrome("https://www.photopea.com/"),
	"web tool pdfescape":               OpenUrlWithChrome("https://www.pdfescape.com/"),
	"web tool smallpdf":                OpenUrlWithChrome("https://smallpdf.com/"),
	"web tool wechat file help":        OpenUrlWithChrome("https://filehelper.weixin.qq.com/"),
	"web video bilibili":               OpenUrlWithChrome("https://www.bilibili.com"),
	"web video cctv5":                  OpenUrlWithChrome("https://tv.cctv.com/live/cctv5"),
	"web video iqiyi":                  OpenUrlWithChrome("https://www.iqiyi.com/"),
	"web video mg tv":                  OpenUrlWithChrome("https://www.mgtv.com/"),
	"web video qq tv":                  OpenUrlWithChrome("https://v.qq.com/"),
	"web video youku":                  OpenUrlWithChrome("https://www.youku.com/"),
	"web video youtube":                OpenUrlWithChrome("https://www.youtube.com"),
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
		OpenUrlWithChrome(content)()
		return
	default:
		SearchFromWeb(content)
	}
}
