package filebrowser

import (
	"embed"
	"errors"
	"net/http"
	"path/filepath"

	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

//go:embed ui.html
var uiFS embed.FS

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.RouterGroup) {
	r.GET("/list", h.List)
	r.GET("/download", h.Download)
	r.GET("/raw", h.Raw)
	r.GET("/ui", h.UI)
	r.GET("", h.UI)
}

func mapErr(c *gin.Context, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, ErrOutsideRoot), errors.Is(err, ErrDenied):
		response.Err(c, response.CodeForbidden, err.Error())
	case errors.Is(err, ErrNotFound):
		response.Err(c, response.CodeNotFound, err.Error())
	case errors.Is(err, ErrTooLarge):
		response.Err(c, response.CodeInvalidParams, err.Error())
	default:
		response.ErrWithInternal(c, response.CodeServerError, "filebrowser error", err)
	}
	return true
}

// List godoc
// @Summary 列出目录
// @Tags 文件浏览
// @Param path query string false "相对 root 的路径(空=root)"
// @Success 200 {object} response.Response
// @Router /api/files/list [get]
func (h *Handler) List(c *gin.Context) {
	res, err := h.svc.ListDir(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	response.Ok(c, res)
}

// Download godoc
// @Summary 下载文件
// @Tags 文件浏览
// @Param path query string true "相对 root 的文件路径"
// @Router /api/files/download [get]
func (h *Handler) Download(c *gin.Context) {
	abs, info, err := h.svc.ResolveFile(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	c.FileAttachment(abs, info.Name())
}

// Raw godoc
// @Summary 在线预览文件
// @Tags 文件浏览
// @Param path query string true "相对 root 的文件路径"
// @Router /api/files/raw [get]
func (h *Handler) Raw(c *gin.Context) {
	abs, info, err := h.svc.ResolveRaw(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	_ = info
	ext := filepath.Ext(abs)
	switch ext {
	case ".md", ".log", ".txt", ".go", ".py", ".sh", ".js", ".ts", ".json",
		".yaml", ".yml", ".toml", ".ini", ".conf", ".html", ".css", ".xml", ".sql":
		c.Header("Content-Type", "text/plain; charset=utf-8")
	}
	c.File(abs)
}

// UI godoc
// @Summary 文件浏览前端
// @Tags 文件浏览
// @Router /api/files/ui [get]
func (h *Handler) UI(c *gin.Context) {
	data, err := uiFS.ReadFile("ui.html")
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "ui not found", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}
