package filebrowser

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"assistant/pkg/cache"
	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

var thumbCache = cache.NewCacheWithMax(8*time.Minute, time.Minute, 500)

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
	r.POST("/download-tgz", h.DownloadTarGz)
	r.GET("/download-tgz", h.DownloadTarGz)
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
	if c.Query("thumb") == "1" {
		h.thumb(c)
		return
	}
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

// thumb 生成缩略图：最长边 ≤ 300px，JPEG 质量 85。
// 结果缓存在内存中(LRU+TTL 10min，最多 500 张)。
func (h *Handler) thumb(c *gin.Context) {
	abs, _, err := h.svc.ResolveFile(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	stat, err := os.Stat(abs)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "stat failed", err)
		return
	}

	key := abs + "@" + stat.ModTime().Format(time.RFC3339Nano) + "@q85"
	if data, ok := thumbCache.Get(key); ok {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Data(http.StatusOK, "image/jpeg", data.([]byte))
		return
	}

	f, err := os.Open(abs)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "open failed", err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		f.Seek(0, 0)
		img, err = webp.Decode(f)
		if err != nil {
			f.Seek(0, 0)
			buf := make([]byte, 512)
			n, _ := io.ReadFull(f, buf)
			ct := http.DetectContentType(buf[:n])
			f.Seek(0, 0)
			c.Header("Content-Type", ct)
			http.ServeContent(c.Writer, c.Request, abs, stat.ModTime(), f)
			return
		}
	}

	bounds := img.Bounds()
	maxSize := 300
	iw, ih := bounds.Dx(), bounds.Dy()
	if iw > ih {
		if iw > maxSize {
			ih = ih * maxSize / iw
			iw = maxSize
		}
	} else {
		if ih > maxSize {
			iw = iw * maxSize / ih
			ih = maxSize
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, iw, ih))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85}); err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "thumbnail encode failed", err)
		return
	}
	data := buf.Bytes()

	thumbCache.Set(key, data)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, "image/jpeg", data)
}

// DownloadTarGz godoc
// @Summary 下载选中文件/目录为 tar.gz
// @Description 将指定路径列表(文件或目录)打包为 tar.gz 下载
// @Tags 文件浏览
// @Accept json
// @Param paths query string false "路径(可重复 path=a&path=b)"
// @Success 200 {file} binary
// @Router /api/files/download-tgz [post]
func (h *Handler) DownloadTarGz(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths"`
	}
	if c.Request.Method == "GET" {
		req.Paths = c.QueryArray("path")
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Err(c, response.CodeInvalidParams, err.Error())
			return
		}
	}
	if len(req.Paths) == 0 {
		response.Err(c, response.CodeInvalidParams, "paths is required")
		return
	}
	tgzPath, err := h.svc.CreateTarGz(req.Paths)
	if mapErr(c, err) {
		return
	}
	defer os.Remove(tgzPath)
	fname := tgzFilename(req.Paths)
	data, err := os.ReadFile(tgzPath)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "read failed", err)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fname))
	c.Data(http.StatusOK, "application/gzip", data)
}

func tgzFilename(paths []string) string {
	if len(paths) == 1 {
		base := filepath.Base(paths[0])
		if base == "" || base == "." {
			return "download.tar.gz"
		}
		return base + ".tar.gz"
	}
	// 多个文件：找共同父目录名
	parent := filepath.Dir(paths[0])
	for _, p := range paths[1:] {
		if filepath.Dir(p) != parent {
			return "download.tar.gz"
		}
	}
	base := filepath.Base(parent)
	if base == "" || base == "." || base == "/" {
		return "download.tar.gz"
	}
	return base + ".tar.gz"
}

// Upload godoc
// @Summary 上传文件(不覆盖，支持多文件)
// @Description multipart/form-data 上传; 字段 file 为文件(可多个), path 为目标目录(相对 root, 空=root)
// @Tags 文件浏览
// @Accept multipart/form-data
// @Param path formData string false "目标目录(相对 root)"
// @Param file formData file true "文件(可多个)"
// @Success 200 {object} response.Response
// @Router /api/files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadBytes)

	form, err := c.MultipartForm()
	if err != nil {
		response.Err(c, response.CodeInvalidParams, "parse form failed: "+err.Error())
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		response.Err(c, response.CodeInvalidParams, "no file provided")
		return
	}

	dir := c.Query("path")
	if dir == "" {
		dir = c.PostForm("path")
	}

	results := make([]gin.H, 0, len(files))
	for _, fh := range files {
		r := h.uploadOne(dir, fh)
		results = append(results, r)
	}
	response.Ok(c, results)
}

func (h *Handler) uploadOne(dir string, fh *multipart.FileHeader) gin.H {
	if fh.Size > MaxUploadBytes {
		return gin.H{"name": fh.Filename, "error": "file too large"}
	}
	name := filepath.Base(fh.Filename)
	dst, finalRel, err := h.svc.CreateUploadFile(dir, name)
	if err != nil {
		return gin.H{"name": fh.Filename, "error": err.Error()}
	}
	src, err := fh.Open()
	if err != nil {
		_ = dst.Close()
		_ = os.Remove(dst.Name())
		return gin.H{"name": fh.Filename, "error": err.Error()}
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		_ = os.Remove(dst.Name())
		return gin.H{"name": fh.Filename, "error": err.Error()}
	}
	if err := dst.Close(); err != nil {
		_ = os.Remove(dst.Name())
		return gin.H{"name": fh.Filename, "error": err.Error()}
	}
	return gin.H{"name": fh.Filename, "path": finalRel, "size": fh.Size}
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
