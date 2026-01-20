package plugins

import (
	"sort"

	"cmds/utils"
)

var ActionMap = map[string]func(){
	"handle copied":                    HandleCopied,
	"format json":                      FormatJson,
	"format sql":                       FormatSql,
	"format yaml":                      FormatYaml,
	"format go":                        FormatGo,
	"get current datetime":             GetCurrentDatetime,
	"get current unix sec":             GetCurrentUnixSec,
	"get host name":                    GetHostName,
	"get ip address":                   GetIPAddress,
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
	"sys display":                      SysDisplay,
	"sys shortcuts":                    SysShortcuts,
	"sys toggle keyboard light":        SysToggleKeyboardLight,
	"sys wifi connect":                 SysWifiConnect,
	"transform datetime to unix sec":   TransformDatetime2UnixSec,
	"transform unix sec to datetime":   TransformUnixSec2DateTime,
	"launch inkscape":                  LaunchInkscape,
	"launch krita":                     LaunchKrita,
	"launch obsidian":                  LaunchObsidian,
	"launch sublime":                   LaunchSublime,
	"launch xournal":                   LaunchXournal,
	"toggle calendar scheduling today": ToggleCalendarSchedulingToday,
	"toggle calendar scheduling":       ToggleCalendarScheduling,
	"toggle calendar":                  ToggleCalendar,
	"toggle clipmenu":                  ToggleClipmenu,
	"toggle flameshot":                 ToggleFlameshot,
	"toggle irssi":                     ToggleIrssi,
	"toggle julia":                     ToggleJulia,
	"toggle lazydocker":                ToggleLazyDocker,
	"toggle music":                     ToggleMusic,
	"toggle netease cloud music":       ToggleNeteaseCloudMusic,
	"toggle passmenu":                  TogglePassmenu,
	"toggle python":                    TogglePython,
	"toggle rec-audio":                 ToggleRecAudio,
	"toggle rec-screen":                ToggleRecScreen,
	"toggle rec-webcam":                ToggleRecWebcam,
	"toggle rec-show":                  ToggleRecShow,
	"toggle screenkey":                 ToggleScreenKey,
	"toggle top":                       ToggleTop,
	"toggle tty_clock":                 ToggleTTYClock,
	"toggle yazi":                      ToggleYazi,
	"web ai chatgpt":                   OpenUrlWithQutebrowser("https://chatgpt.com/"),
	"web ai doubao":                    OpenUrlWithQutebrowser("https://www.doubao.com/"),
	"web ai gemini":                    OpenUrlWithQutebrowser("https://gemini.google.com/app"),
	"web cs color":                     OpenUrlWithQutebrowser("https://coolors.co/palettes/trending"),
	"web cs github gist":               OpenUrlWithQutebrowser("https://gist.github.com/"),
	"web cs github repos":              OpenUrlWithQutebrowser("https://github.com/zetatez?tab=repositories"),
	"web cs github":                    OpenUrlWithQutebrowser("https://github.com"),
	"web cs leetcode":                  OpenUrlWithQutebrowser("https://leetcode.cn/problemset/"),
	"web file browser":                 OpenUrlWithQutebrowser("http://127.0.0.1:5080/files"),
	"web map":                          OpenUrlWithQutebrowser("https://ditu.amap.com/"),
	"web mirror tsinghua":              OpenUrlWithQutebrowser("https://mirrors.tuna.tsinghua.edu.cn/"),
	"web mirror sjtu":                  OpenUrlWithQutebrowser("https://ftp.sjtu.edu.cn/"),
	"web scholar arxiv":                OpenUrlWithQutebrowser("https://arxiv.org"),
	"web scholar scholar":              OpenUrlWithQutebrowser("https://scholar.google.com"),
	"web scholar wolframalpha":         OpenUrlWithQutebrowser("https://www.wolframalpha.com"),
	"web social douban":                OpenUrlWithQutebrowser("https://www.douban.com/"),
	"web social instagram":             OpenUrlWithQutebrowser("https://www.instagram.com/explore/"),
	"web social wechat file help":      OpenUrlWithQutebrowser("https://filehelper.weixin.qq.com/"),
	"web social wechat":                OpenUrlWithQutebrowser("https://web.wechat.com/"),
	"web suckless":                     OpenUrlWithQutebrowser("https://dwm.suckless.org/"),
	"web translate auto to en":         OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=en"),
	"web translate auto to zh":         OpenUrlWithQutebrowser("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web video bilibili":               OpenUrlWithQutebrowser("https://www.bilibili.com"),
	"web video cctv5":                  OpenUrlWithQutebrowser("https://tv.cctv.com/live/cctv5"),
	"web video youtube":                OpenUrlWithQutebrowser("https://www.youtube.com"),
	"web vpn shadowsocks":              OpenUrlWithQutebrowser("https://portal.shadowsocks.nz/login"),
	"web wiki archlinux":               OpenUrlWithQutebrowser("https://wiki.archlinux.org"),
	"web google calendar":              OpenUrlWithQutebrowser("https://calendar.google.com/calendar/u/0/r/month/2026/1/1?pli=1"),
	"web google mail":                  OpenUrlWithQutebrowser("https://accounts.google.com/b/0/AddMailService"),
	"web google map":                   OpenUrlWithQutebrowser("https://www.google.com/maps/place/shanghai"),
	"web google":                       OpenUrlWithQutebrowser("https://www.google.com/"),
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
