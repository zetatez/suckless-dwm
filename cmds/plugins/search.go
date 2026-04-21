package plugins

import (
	"sort"

	"cmds/utils"
)

var ActionMap = map[string]func(){
	"handle clipboard":                 HandleClipboard,
	"format json":                      FormatJson,
	"format sql":                       FormatSql,
	"format yaml":                      FormatYaml,
	"format go":                        FormatGo,
	"get cur datetime":                 GetCurDatetime,
	"get cur unix sec":                 GetCurUnixSec,
	"get ip address":                   GetIPAddress,
	"send clipboard to feishu robot":   SendClipboardToFeishuRobot,
	"note script":                      NoteScripts,
	"note todo":                        NoteToDo,
	"note monthly work":                NoteMonthlyWork,
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
	"conversion datetime to unix sec":  ConversionDatetime2UnixSec,
	"conversion unix sec to datetime":  ConversionUnixSec2DateTime,
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
	"web ai chatgpt":                   OpenUrlWithChrome("https://chatgpt.com/"),
	"web ai doubao":                    OpenUrlWithChrome("https://www.doubao.com/"),
	"web ai gemini":                    OpenUrlWithChrome("https://gemini.google.com/app"),
	"web cs color":                     OpenUrlWithChrome("https://coolors.co/palettes/trending"),
	"web cs github gist":               OpenUrlWithChrome("https://gist.github.com/"),
	"web cs github repos":              OpenUrlWithChrome("https://github.com/zetatez?tab=repositories"),
	"web cs github":                    OpenUrlWithChrome("https://github.com"),
	"web cs leetcode":                  OpenUrlWithChrome("https://leetcode.cn/problemset/"),
	"web file browser":                 OpenUrlWithChrome("http://127.0.0.1:5080/files"),
	"web map":                          OpenUrlWithChrome("https://ditu.amap.com/"),
	"web mirror tsinghua":              OpenUrlWithChrome("https://mirrors.tuna.tsinghua.edu.cn/"),
	"web mirror sjtu":                  OpenUrlWithChrome("https://ftp.sjtu.edu.cn/"),
	"web scholar arxiv":                OpenUrlWithChrome("https://arxiv.org"),
	"web scholar scholar":              OpenUrlWithChrome("https://scholar.google.com"),
	"web scholar wolframalpha":         OpenUrlWithChrome("https://www.wolframalpha.com"),
	"web social douban":                OpenUrlWithChrome("https://www.douban.com/"),
	"web social instagram":             OpenUrlWithChrome("https://www.instagram.com/explore/"),
	"web social wechat file help":      OpenUrlWithChrome("https://filehelper.weixin.qq.com/"),
	"web social wechat":                OpenUrlWithChrome("https://web.wechat.com/"),
	"web suckless":                     OpenUrlWithChrome("https://dwm.suckless.org/"),
	"web translate to en":              OpenUrlWithChrome("https://translate.google.com/?sl=auto&tl=en"),
	"web translate to zh":              OpenUrlWithChrome("https://translate.google.com/?sl=auto&tl=zh-CN"),
	"web video bilibili":               OpenUrlWithChrome("https://www.bilibili.com"),
	"web video cctv5":                  OpenUrlWithChrome("https://tv.cctv.com/live/cctv5"),
	"web video youtube":                OpenUrlWithChrome("https://www.youtube.com"),
	"web vpn shadowsocks":              OpenUrlWithChrome("https://portal.shadowsocks.nz/login"),
	"web wiki archlinux":               OpenUrlWithChrome("https://wiki.archlinux.org"),
	"web google calendar":              OpenUrlWithChrome("https://calendar.google.com/calendar/u/0/r/month/2026/1/1?pli=1"),
	"web google mail":                  OpenUrlWithChrome("https://accounts.google.com/b/0/AddMailService"),
	"web google map":                   OpenUrlWithChrome("https://www.google.com/maps/place/shanghai"),
	"web google":                       OpenUrlWithChrome("https://www.google.com/"),
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
