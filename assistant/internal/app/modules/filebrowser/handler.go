package filebrowser

import (
	"embed"
	"errors"
	"io"
	"net/http"
	"os"
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
	r.POST("/upload", h.Upload)
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
	case errors.Is(err, ErrTooLarge), errors.Is(err, ErrBadName):
		response.Err(c, response.CodeInvalidParams, err.Error())
	case errors.Is(err, ErrExists):
		response.Err(c, response.CodeForbidden, err.Error())
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
	abs, _, err := h.svc.ResolveRaw(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	f, err := os.Open(abs)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "open failed", err)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	buf := make([]byte, 512)
	n, _ := io.ReadFull(f, buf)
	ct := http.DetectContentType(buf[:n])
	f.Seek(0, 0)

	c.Header("Content-Type", ct)
	http.ServeContent(c.Writer, c.Request, abs, stat.ModTime(), f)
}

// Upload godoc
// @Summary 上传文件(不覆盖)
// @Description multipart/form-data 上传; 字段 file 为文件, path 为目标目录(相对 root, 空=root)
// @Tags 文件浏览
// @Accept multipart/form-data
// @Param path formData string false "目标目录(相对 root)"
// @Param file formData file true "文件"
// @Success 200 {object} response.Response
// @Router /api/files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadBytes)

	fh, err := c.FormFile("file")
	if err != nil {
		response.Err(c, response.CodeInvalidParams, "file is required: "+err.Error())
		return
	}
	if fh.Size > MaxUploadBytes {
		response.Err(c, response.CodeInvalidParams, "file too large")
		return
	}

	dir := c.Query("path")
	if dir == "" {
		dir = c.PostForm("path")
	}
	name := filepath.Base(fh.Filename)

	dst, finalRel, err := h.svc.CreateUploadFile(dir, name)
	if mapErr(c, err) {
		return
	}
	src, err := fh.Open()
	if err != nil {
		_ = dst.Close()
		_ = os.Remove(dst.Name())
		response.ErrWithInternal(c, response.CodeServerError, "open upload failed", err)
		return
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		_ = os.Remove(dst.Name())
		response.ErrWithInternal(c, response.CodeServerError, "write failed", err)
		return
	}
	if err := dst.Close(); err != nil {
		_ = os.Remove(dst.Name())
		response.ErrWithInternal(c, response.CodeServerError, "close failed", err)
		return
	}
	response.Ok(c, gin.H{"path": finalRel, "size": fh.Size})
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
