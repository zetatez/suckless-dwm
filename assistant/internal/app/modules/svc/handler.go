package svc

import (
	"fmt"

	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) trigger(c *gin.Context, name string, fn func() error) {
	if err := fn(); err != nil {
		h.svc.notify(fmt.Sprintf("%s failed: %v", name, err))
		response.ErrWithInternal(c, response.CodeServerError, name+" failed", err)
		return
	}
	response.Ok(c, gin.H{"status": "done"})
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/sys-shortcut", h.SysShortcut)
	r.POST("/format", h.Format)
	r.POST("/note", h.Note)
	r.GET("/get-datetime", h.GetDatetime)
	r.GET("/get-cur-unix-sec", h.GetCurUnixSec)
	r.POST("/convert-datetime", h.ConvertDatetime)
	r.GET("/get-ip", h.GetIP)
	r.POST("/feishu-send", h.FeishuSend)
	r.POST("/toggle", h.Toggle)
	r.POST("/launch", h.Launch)
	r.POST("/search-web", h.SearchWeb)
	r.POST("/search-books-online", h.SearchBooksOnline)
	r.POST("/search-videos-online", h.SearchVideosOnline)
	r.POST("/sys-bluetooth-connect", h.SysBluetoothConnect)
	r.POST("/sys-bluetooth-scan-connect", h.SysBluetoothScanConnect)
	r.POST("/sys-bluetooth-disconnect", h.SysBluetoothDisconnect)
	r.POST("/sys-wifi-connect", h.SysWifiConnect)
	r.POST("/sys-display", h.SysDisplay)
	r.POST("/sys-keyboard-light", h.SysKeyboardLight)
	r.POST("/open-url", h.OpenURL)
	r.POST("/open-url-as-app", h.OpenURLAsApp)
	r.POST("/sys-ssh-connect", h.SysSSHConnect)
	r.POST("/handle-clipboard", h.HandleClipboard)
	r.POST("/sys-volume-up", h.SysVolumeUp)
	r.POST("/sys-volume-down", h.SysVolumeDown)
	r.POST("/sys-volume-toggle", h.SysVolumeToggle)
	r.POST("/sys-micro-up", h.SysMicroUp)
	r.POST("/sys-micro-down", h.SysMicroDown)
	r.POST("/sys-micro-toggle", h.SysMicroToggle)
	r.POST("/sys-display-light-up", h.SysDisplayLightUp)
	r.POST("/sys-display-light-down", h.SysDisplayLightDown)
	r.POST("/sys-reset", h.SysReset)
	r.POST("/sys-kill", h.SysKill)
	r.POST("/sys-open-terminal", h.SysOpenTerminal)
	r.POST("/file-search", h.FileSearch)
	r.POST("/file-search-content", h.FileSearchContent)
	r.POST("/file-search-book", h.FileSearchBook)
	r.POST("/file-search-media", h.FileSearchMedia)
	r.POST("/file-search-wiki", h.FileSearchWiki)
	r.POST("/file-open-images", h.FileOpenImages)
	r.POST("/snip-fzf", h.SnipFzf)
	r.POST("/snip-create", h.SnipCreate)
	r.POST("/search", h.Search)
}

// Power godoc
// @Summary 电源操作菜单
// @Description 触发 rofi 电源菜单(suspend/poweroff/reboot/off-display/slock)
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-shortcut [post]
func (h *Handler) SysShortcut(c *gin.Context) { h.trigger(c, "power", h.svc.SysShortcut) }

// Format godoc
// @Summary 格式化代码
// @Description 从剪贴板读取代码，格式化后写回剪贴板(json/yaml/sql/go)
// @Tags 工具
// @Accept json
// @Param body body FormatRequest true "语言类型"
// @Success 200 {object} response.Response
// @Router /api/svr/format [post]
func (h *Handler) Format(c *gin.Context) {
	var req FormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "language is required")
		return
	}
	result, err := h.svc.Format(req.Language)
	if err != nil {
		h.svc.notify(fmt.Sprintf("format failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "format failed", err)
		return
	}
	response.Ok(c, gin.H{"language": req.Language, "result": result})
}

// Note godoc
// @Summary 笔记
// @Description 写入 header 后用 nvim 打开编辑
// @Tags 工具
// @Accept json
// @Param body body NoteRequest true "笔记类型(todo/scripts/monthly_work)"
// @Success 200 {object} response.Response
// @Router /api/svr/note [post]
func (h *Handler) Note(c *gin.Context) {
	var req NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "type is required")
		return
	}
	if err := h.svc.Note(req.Type); err != nil {
		h.svc.notify(fmt.Sprintf("note failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "note failed", err)
		return
	}
	response.Ok(c, gin.H{"type": req.Type, "status": "appended"})
}

// GetDatetime godoc
// @Summary 获取当前时间
// @Description 获取当前日期时间并写入剪贴板
// @Tags 工具
// @Success 200 {object} response.Response
// @Router /api/svr/get-datetime [get]
func (h *Handler) GetDatetime(c *gin.Context) {
	response.Ok(c, h.svc.GetDatetime())
}

// GetCurUnixSec godoc
// @Summary 获取当前 Unix 时间戳
// @Description 获取当前 Unix 时间戳并写入剪贴板
// @Tags 工具
// @Success 200 {object} response.Response
// @Router /api/svr/get-cur-unix-sec [get]
func (h *Handler) GetCurUnixSec(c *gin.Context) {
	response.Ok(c, gin.H{"unix": h.svc.GetCurUnixSec()})
}

// ConvertDatetime godoc
// @Summary 时间格式转换
// @Description 从剪贴板读取时间字符串，转换后写回剪贴板
// @Tags 工具
// @Accept json
// @Param body body DatetimeConvertRequest true "转换格式"
// @Success 200 {object} response.Response
// @Router /api/svr/convert-datetime [post]
func (h *Handler) ConvertDatetime(c *gin.Context) {
	var req DatetimeConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "from and to are required")
		return
	}
	result, err := h.svc.ConvertDatetime(req.From, req.To)
	if err != nil {
		h.svc.notify(fmt.Sprintf("convert failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "convert failed", err)
		return
	}
	response.Ok(c, gin.H{"result": result})
}

// GetIP godoc
// @Summary 获取 IP 地址
// @Description 获取网卡 IP 地址并写入剪贴板
// @Tags 网络
// @Param interface query string false "网卡名称(默认 wlan0)"
// @Success 200 {object} response.Response
// @Router /api/svr/get-ip [get]
func (h *Handler) GetIP(c *gin.Context) {
	iface := c.DefaultQuery("interface", "")
	ips, err := h.svc.GetIP(iface)
	if err != nil {
		h.svc.notify(fmt.Sprintf("get ip failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "get IP failed", err)
		return
	}
	response.Ok(c, gin.H{"interface": iface, "ips": ips})
}

// FeishuSend godoc
// @Summary 发送飞书消息
// @Description 从剪贴板读取内容并发送到飞书机器人
// @Tags 通信
// @Success 200 {object} response.Response
// @Router /api/svr/feishu-send [post]
func (h *Handler) FeishuSend(c *gin.Context) { h.trigger(c, "feishu", h.svc.FeishuSend) }

// Toggle godoc
// @Summary 切换进程
// @Description 如果进程运行则杀死，否则启动
// @Tags 进程
// @Accept json
// @Param body body ToggleRequest true "进程名"
// @Success 200 {object} response.Response
// @Router /api/svr/toggle [post]
func (h *Handler) Toggle(c *gin.Context) {
	var req ToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "process is required")
		return
	}
	response.Ok(c, gin.H{"process": req.Process, "status": h.svc.Toggle(req.Process)})
}

// Launch godoc
// @Summary 启动应用
// @Description 启动指定应用
// @Tags 进程
// @Accept json
// @Param body body LaunchRequest true "启动命令"
// @Success 200 {object} response.Response
// @Router /api/svr/launch [post]
func (h *Handler) Launch(c *gin.Context) {
	var req LaunchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "command is required")
		return
	}
	if err := h.svc.Launch(req.Command); err != nil {
		h.svc.notify(fmt.Sprintf("launch failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "launch failed", err)
		return
	}
	response.Ok(c, gin.H{"command": req.Command, "status": "launched"})
}

// SearchWeb godoc
// @Summary 网页搜索
// @Description 用 Chrome 打开 Google 搜索结果
// @Tags 搜索
// @Accept json
// @Param body body object true "搜索关键词"
// @Success 200 {object} response.Response
// @Router /api/svr/search-web [post]
func (h *Handler) SearchWeb(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "query is required")
		return
	}
	if err := h.svc.SearchWeb(req.Query); err != nil {
		h.svc.notify(fmt.Sprintf("search failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "search failed", err)
		return
	}
	response.Ok(c, gin.H{"query": req.Query, "status": "opened"})
}

// SearchBooksOnline godoc
// @Summary 图书搜索
// @Description 同时打开 openlibrary 和 z-lib 的图书搜索结果
// @Tags 搜索
// @Accept json
// @Param body body object true "搜索关键词"
// @Success 200 {object} response.Response
// @Router /api/svr/search-books-online [post]
func (h *Handler) SearchBooksOnline(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "query is required")
		return
	}
	if err := h.svc.SearchBooksOnline(req.Query); err != nil {
		h.svc.notify(fmt.Sprintf("search books failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "search books failed", err)
		return
	}
	response.Ok(c, gin.H{"query": req.Query, "status": "opened"})
}

// SearchVideosOnline godoc
// @Summary 视频搜索
// @Description 同时打开 Bilibili 和 YouTube 的视频搜索结果
// @Tags 搜索
// @Accept json
// @Param body body object true "搜索关键词"
// @Success 200 {object} response.Response
// @Router /api/svr/search-videos-online [post]
func (h *Handler) SearchVideosOnline(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "query is required")
		return
	}
	if err := h.svc.SearchVideosOnline(req.Query); err != nil {
		h.svc.notify(fmt.Sprintf("search videos failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "search videos failed", err)
		return
	}
	response.Ok(c, gin.H{"query": req.Query, "status": "opened"})
}

// BluetoothConnect godoc
// @Summary 蓝牙连接
// @Description 触发 rofi 选择已配对设备并连接
// @Tags 蓝牙
// @Success 200 {object} response.Response
// @Router /api/svr/sys-bluetooth-connect [post]
func (h *Handler) SysBluetoothConnect(c *gin.Context) {
	h.trigger(c, "bluetooth", h.svc.SysBluetoothConnect)
}

// BluetoothScanConnect godoc
// @Summary 蓝牙扫描并连接
// @Description 扫描附近设备后 rofi 选择并连接
// @Tags 蓝牙
// @Success 200 {object} response.Response
// @Router /api/svr/sys-bluetooth-scan-connect [post]
func (h *Handler) SysBluetoothScanConnect(c *gin.Context) {
	h.trigger(c, "bluetooth", h.svc.SysBluetoothScanConnect)
}

// BluetoothDisconnect godoc
// @Summary 蓝牙断开
// @Description 触发 rofi 选择已连接设备并断开
// @Tags 蓝牙
// @Success 200 {object} response.Response
// @Router /api/svr/sys-bluetooth-disconnect [post]
func (h *Handler) SysBluetoothDisconnect(c *gin.Context) {
	h.trigger(c, "bluetooth", h.svc.SysBluetoothDisconnect)
}

// WifiConnect godoc
// @Summary WiFi 连接
// @Description 触发 rofi 选择 WiFi 并输入密码后连接
// @Tags 网络
// @Success 200 {object} response.Response
// @Router /api/svr/sys-wifi-connect [post]
func (h *Handler) SysWifiConnect(c *gin.Context) { h.trigger(c, "wifi", h.svc.SysWifiConnect) }

// Display godoc
// @Summary 显示器布局
// @Description 触发 rofi 选择显示器布局
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-display [post]
func (h *Handler) SysDisplay(c *gin.Context) { h.trigger(c, "display", h.svc.SysDisplay) }

// KeyboardLight godoc
// @Summary 键盘背光
// @Description 切换 ThinkPad 键盘背光
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-keyboard-light [post]
func (h *Handler) SysKeyboardLight(c *gin.Context) {
	val, err := h.svc.SysKeyboardLight()
	if err != nil {
		h.svc.notify(fmt.Sprintf("keyboard light failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "toggle keyboard light failed", err)
		return
	}
	response.Ok(c, gin.H{"brightness": val})
}

// OpenURL godoc
// @Summary 打开 URL
// @Description 用指定浏览器打开 URL(支持 chrome/qutebrowser)
// @Tags 网络
// @Accept json
// @Param body body object true "浏览器和URL"
// @Success 200 {object} response.Response
// @Router /api/svr/open-url [post]
func (h *Handler) OpenURL(c *gin.Context) {
	var req struct {
		Browser string `json:"browser"`
		URL     string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "url is required")
		return
	}
	browser := req.Browser
	if browser == "" {
		browser = "chrome"
	}
	if err := h.svc.OpenURL(browser, req.URL); err != nil {
		h.svc.notify(fmt.Sprintf("open url failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "open url failed", err)
		return
	}
	response.Ok(c, gin.H{"browser": browser, "url": req.URL, "status": "opened"})
}

// OpenURLAsApp godoc
// @Summary 应用模式打开 URL
// @Description 用浏览器应用模式打开 URL(chrome-app/qutebrowser)
// @Tags 网络
// @Accept json
// @Param body body object true "浏览器和URL"
// @Success 200 {object} response.Response
// @Router /api/svr/open-url-as-app [post]
func (h *Handler) OpenURLAsApp(c *gin.Context) {
	var req struct {
		Browser string `json:"browser"`
		URL     string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, "url is required")
		return
	}
	if req.Browser == "" {
		req.Browser = "chrome"
	}
	if err := h.svc.OpenURLAsApp(req.Browser, req.URL); err != nil {
		h.svc.notify(fmt.Sprintf("open url as app failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "open url as app failed", err)
		return
	}
	response.Ok(c, gin.H{"browser": req.Browser, "url": req.URL, "status": "opened"})
}

// SSHConnect godoc
// @Summary SSH 连接
// @Description 触发 rofi 选择 SSH 主机并输入密码后连接
// @Tags 网络
// @Success 200 {object} response.Response
// @Router /api/svr/sys-ssh-connect [post]
func (h *Handler) SysSSHConnect(c *gin.Context) { h.trigger(c, "ssh", h.svc.SysSSHConnect) }

// HandleClipboard godoc
// @Summary 智能剪贴板
// @Description 从剪贴板读取内容，自动判断: 文件路径、URL 或搜索
// @Tags 工具
// @Success 200 {object} response.Response
// @Router /api/svr/handle-clipboard [post]
func (h *Handler) HandleClipboard(c *gin.Context) {
	action, err := h.svc.HandleClipboard()
	if err != nil {
		h.svc.notify(fmt.Sprintf("clipboard handle failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "handle clipboard failed", err)
		return
	}
	response.Ok(c, gin.H{"action": action})
}

// VolumeUp godoc
// @Summary 音量+
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-volume-up [post]
func (h *Handler) SysVolumeUp(c *gin.Context) { h.trigger(c, "volume up", h.svc.SysVolumeUp) }

// VolumeDown godoc
// @Summary 音量-
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-volume-down [post]
func (h *Handler) SysVolumeDown(c *gin.Context) { h.trigger(c, "volume down", h.svc.SysVolumeDown) }

// VolumeToggle godoc
// @Summary 静音切换
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-volume-toggle [post]
func (h *Handler) SysVolumeToggle(c *gin.Context) {
	h.trigger(c, "volume toggle", h.svc.SysVolumeToggle)
}

// MicroUp godoc
// @Summary 麦克风+
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-micro-up [post]
func (h *Handler) SysMicroUp(c *gin.Context) { h.trigger(c, "micro up", h.svc.SysMicroUp) }

// MicroDown godoc
// @Summary 麦克风-
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-micro-down [post]
func (h *Handler) SysMicroDown(c *gin.Context) { h.trigger(c, "micro down", h.svc.SysMicroDown) }

// MicroToggle godoc
// @Summary 麦克风开关
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-micro-toggle [post]
func (h *Handler) SysMicroToggle(c *gin.Context) { h.trigger(c, "micro toggle", h.svc.SysMicroToggle) }

// DisplayLightUp godoc
// @Summary 屏幕亮度+
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-display-light-up [post]
func (h *Handler) SysDisplayLightUp(c *gin.Context) {
	h.trigger(c, "display light up", h.svc.SysDisplayLightUp)
}

// DisplayLightDown godoc
// @Summary 屏幕亮度-
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-display-light-down [post]
func (h *Handler) SysDisplayLightDown(c *gin.Context) {
	h.trigger(c, "display light down", h.svc.SysDisplayLightDown)
}

// SysReset godoc
// @Summary 系统重置
// @Description 重置亮度/音量/麦克风/键盘速率到默认值
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-reset [post]
func (h *Handler) SysReset(c *gin.Context) { h.trigger(c, "sys reset", h.svc.SysReset) }

// SysKill godoc
// @Summary 进程管理器
// @Description 触发 fzf 交互式进程管理和杀死
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-kill [post]
func (h *Handler) SysKill(c *gin.Context) { h.trigger(c, "sys kill", h.svc.SysKill) }

// SysOpenTerminal godoc
// @Summary 打开终端
// @Description 触发 fzf 选择目录并打开终端
// @Tags 系统
// @Success 200 {object} response.Response
// @Router /api/svr/sys-open-terminal [post]
func (h *Handler) SysOpenTerminal(c *gin.Context) {
	h.trigger(c, "open terminal", h.svc.SysOpenTerminal)
}

// FileSearch godoc
// @Summary 文件搜索
// @Description 触发 fzf 全盘文件搜索并打开
// @Tags 文件
// @Success 200 {object} response.Response
// @Router /api/svr/file-search [post]
func (h *Handler) FileSearch(c *gin.Context) { h.trigger(c, "file search", h.svc.FileSearch) }

// FileSearchContent godoc
// @Summary 文件内容搜索
// @Description 触发 rg+fzf 文件内容搜索
// @Tags 文件
// @Success 200 {object} response.Response
// @Router /api/svr/file-search-content [post]
func (h *Handler) FileSearchContent(c *gin.Context) {
	h.trigger(c, "file search content", h.svc.FileSearchContent)
}

// FileSearchBook godoc
// @Summary 电子书搜索
// @Description 触发 fzf 搜索电子书(pdf/epub/djvu)
// @Tags 文件
// @Success 200 {object} response.Response
// @Router /api/svr/file-search-book [post]
func (h *Handler) FileSearchBook(c *gin.Context) {
	h.trigger(c, "file search book", h.svc.FileSearchBook)
}

// FileSearchMedia godoc
// @Summary 媒体文件搜索
// @Description 触发 fzf 搜索图片/音频/视频文件
// @Tags 文件
// @Success 200 {object} response.Response
// @Router /api/svr/file-search-media [post]
func (h *Handler) FileSearchMedia(c *gin.Context) {
	h.trigger(c, "file search media", h.svc.FileSearchMedia)
}

// FileSearchWiki godoc
// @Summary Wiki 搜索
// @Description 触发 fzf 搜索 Markdown 笔记文件
// @Tags 文件
// @Success 200 {object} response.Response
// @Router /api/svr/file-search-wiki [post]
func (h *Handler) FileSearchWiki(c *gin.Context) {
	h.trigger(c, "file search wiki", h.svc.FileSearchWiki)
}

// FileOpenImages godoc
// @Summary 打开图片
// @Description 用 sxiv 打开目录下所有图片
// @Tags 文件
// @Param dir query string false "目录路径(默认当前目录)"
// @Success 200 {object} response.Response
// @Router /api/svr/file-open-images [post]
func (h *Handler) FileOpenImages(c *gin.Context) {
	dir := c.DefaultQuery("dir", ".")
	if err := h.svc.FileOpenImages(dir); err != nil {
		h.svc.notify(fmt.Sprintf("open images failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "open images failed", err)
		return
	}
	response.Ok(c, gin.H{"status": "done"})
}

// SnipFzf godoc
// @Summary 选择并复制 snip
// @Description fzf 选择 snip 文件, 复制内容到剪贴板
// @Tags 工具
// @Success 200 {object} response.Response
// @Router /api/svr/snip-fzf [post]
func (h *Handler) SnipFzf(c *gin.Context) {
	if err := h.svc.SnipFzf(); err != nil {
		h.svc.notify(fmt.Sprintf("snip fzf failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "snip fzf failed", err)
		return
	}
	response.Ok(c, gin.H{"status": "done"})
}

// SnipCreate godoc
// @Summary 创建 snip
// @Description rofi 输入名称, nvim 编辑内容
// @Tags 工具
// @Param name query string false "snip 名称"
// @Success 200 {object} response.Response
// @Router /api/svr/snip-create [post]
func (h *Handler) SnipCreate(c *gin.Context) {
	name := c.Query("name")
	if err := h.svc.SnipCreate(name); err != nil {
		h.svc.notify(fmt.Sprintf("snip create failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "snip create failed", err)
		return
	}
	response.Ok(c, gin.H{"status": "done"})
}

// Search godoc
// @Summary 搜索并执行动作
// @Description rofi 列出所有可用命令并执行
// @Tags 工具
// @Success 200 {object} response.Response
// @Router /api/svr/search [post]
func (h *Handler) Search(c *gin.Context) {
	if err := h.svc.Search(); err != nil {
		h.svc.notify(fmt.Sprintf("search failed: %v", err))
		response.ErrWithInternal(c, response.CodeServerError, "search failed", err)
		return
	}
	response.Ok(c, gin.H{"status": "done"})
}
