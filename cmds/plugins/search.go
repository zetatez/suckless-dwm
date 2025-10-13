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
	"web fin antfin":                   OpenUrlWithQutebrowser("https://caifu.antfin.com/"),
	"web fin eastmoney fund":           OpenUrlWithQutebrowser("https://fund.eastmoney.com/"),
	"web fin eastmoney":                OpenUrlWithQutebrowser("https://www.eastmoney.com/"),
	"web fin futu":                     OpenUrlWithQutebrowser("https://news.futunn.com/en/main/live?lang=zh-CN"),
	"web fin tonghuashun":              OpenUrlWithQutebrowser("https://www.10jqka.com.cn/"),
	"web fin xueqiu":                   OpenUrlWithQutebrowser("https://xueqiu.com/"),
	"web github gist":                  OpenUrlWithQutebrowser("https://gist.github.com/"),
	"web github repos":                 OpenUrlWithQutebrowser("https://github.com/zetatez?tab=repositories"),
	"web github":                       OpenUrlWithQutebrowser("https://github.com"),
	"web map gaode":                    OpenUrlWithQutebrowser("https://ditu.amap.com/"),
	"web map google":                   OpenUrlWithQutebrowser("https://www.google.com/maps/place/shanghai"),
	"web mirror aliyun":                OpenUrlWithQutebrowser("https://developer.aliyun.com/mirror"),
	"web social douban":                OpenUrlWithQutebrowser("https://www.douban.com/"),
	"web social instagram":             OpenUrlWithQutebrowser("https://www.instagram.com/explore/"),
	"web translate auto to en":         OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=en"),
	"web translate auto to zh":         OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web video bilibili":               OpenUrlWithQutebrowser("https://www.bilibili.com"),
	"web video cctv5":                  OpenUrlWithQutebrowser("https://tv.cctv.com/live/cctv5"),
	"web video youtube":                OpenUrlWithQutebrowser("https://www.youtube.com"),
	"web web study arxiv":              OpenUrlWithQutebrowser("https://arxiv.org"),
	"web web study geekbang":           OpenUrlWithQutebrowser("https://time.geekbang.org/"),
	"web web study leetcode":           OpenUrlWithQutebrowser("https://leetcode.cn/problemset/"),
	"web web study scholar":            OpenUrlWithQutebrowser("https://scholar.google.com"),
	"web web study shanbay":            OpenUrlWithQutebrowser("https://www.shanbay.com/"),
	"web web study wolframalpha":       OpenUrlWithQutebrowser("https://www.wolframalpha.com"),
	"web web study youdao":             OpenUrlWithQutebrowser("https://dict.youdao.com/"),
	"web wechat file help":             OpenUrlWithQutebrowser("https://filehelper.weixin.qq.com/"),
	"web wechat":                       OpenUrlWithQutebrowser("https://web.wechat.com/"),
	"web wiki archlinux":               OpenUrlWithQutebrowser("https://wiki.archlinux.org"),
	"web color":                        OpenUrlWithQutebrowser("https://coolors.co/palettes/trending"),
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
