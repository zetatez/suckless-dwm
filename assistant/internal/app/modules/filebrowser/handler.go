package filebrowser

import (
	"bytes"
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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"assistant/pkg/cache"
	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

var thumbCache = cache.NewCacheWithMax(8*time.Minute, time.Minute, 320)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.RouterGroup) {
	r.GET("", func(c *gin.Context) { renderUI(c.Writer) })
	r.GET("/download", h.Download)
	r.GET("/download-tgz", h.DownloadTarGz)
	r.GET("/list", h.List)
	r.GET("/raw", h.Raw)
	r.GET("/trash", h.ListTrash)
	r.POST("/copy", h.Copy)
	r.POST("/delete", h.Delete)
	r.POST("/mkdir", h.Mkdir)
	r.POST("/move", h.Move)
	r.POST("/rename", h.Rename)
	r.POST("/touch", h.Touch)
	r.POST("/trash/delete", h.PermanentDelete)
	r.POST("/trash/restore", h.RestoreTrash)
	r.POST("/upload", h.Upload)
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
// @Router /api/filebrowser/list [get]
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
// @Success 200 {file} binary
// @Router /api/filebrowser/download [get]
func (h *Handler) Download(c *gin.Context) {
	abs, info, err := h.svc.ResolveFile(c.Query("path"))
	if mapErr(c, err) {
		return
	}
	c.FileAttachment(abs, info.Name())
}

var mimeMap = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".mp3":  "audio/mpeg",
	".flac": "audio/flac",
	".wav":  "audio/wav",
	".ogg":  "audio/ogg",
	".aac":  "audio/aac",
	".m4a":  "audio/mp4",
	".wma":  "audio/x-ms-wma",
	".opus": "audio/opus",
}

// Raw godoc
// @Summary 在线预览文件
// @Tags 文件浏览
// @Param path query string true "相对 root 的文件路径"
// @Param thumb query string false "1=生成缩略图(仅图片)"
// @Success 200 {file} binary
// @Router /api/filebrowser/raw [get]
func (h *Handler) Raw(c *gin.Context) {
	if c.Query("thumb") == "1" {
		h.thumb(c)
		return
	}
	abs, info, err := h.svc.ResolveRaw(c.Query("path"))
	if mapErr(c, err) {
		return
	}

	etag := fmt.Sprintf(`"%d-%d"`, info.ModTime().UnixNano(), info.Size())
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	f, err := os.Open(abs)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "open failed", err)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	ext := strings.ToLower(filepath.Ext(abs))
	ct := mimeMap[ext]
	if ct == "" {
		buf := make([]byte, 512)
		n, _ := io.ReadFull(f, buf)
		f.Seek(0, 0)
		ct = http.DetectContentType(buf[:n])
	}

	c.Header("ETag", etag)
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

	// PDF: render first page via pdftoppm
	if strings.HasSuffix(strings.ToLower(abs), ".pdf") {
		tmpPrefix := filepath.Join(os.TempDir(), "assistant-pdf-"+strconv.FormatInt(time.Now().UnixNano(), 36))
		cmd := exec.Command("pdftoppm", "-f", "1", "-l", "1", "-scale-to", "600", "-png", abs, tmpPrefix)
		out, err := cmd.CombinedOutput()
		if err != nil {
			response.ErrWithInternal(c, response.CodeServerError, "pdftoppm failed: "+string(out), err)
			return
		}
		matches, _ := filepath.Glob(tmpPrefix + "-*.png")
		if len(matches) == 0 {
			response.ErrWithInternal(c, response.CodeServerError, "pdftoppm output not found: "+string(out), nil)
			return
		}
		outPath := matches[0]
		data, err := os.ReadFile(outPath)
		os.Remove(outPath)
		if err != nil {
			response.ErrWithInternal(c, response.CodeServerError, "read pdftoppm output", err)
			return
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			response.ErrWithInternal(c, response.CodeServerError, "decode pdftoppm output", err)
			return
		}
		h.writeThumb(c, key, img)
		return
	}

	// EPUB: render first page via mutool
	if strings.HasSuffix(strings.ToLower(abs), ".epub") {
		cmd := exec.Command("mutool", "draw", "-F", "png", "-o", "-", abs, "1")
		data, err := cmd.Output()
		if err != nil {
			response.ErrWithInternal(c, response.CodeServerError, "mutool failed", err)
			return
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			response.ErrWithInternal(c, response.CodeServerError, "decode mutool output", err)
			return
		}
		h.writeThumb(c, key, img)
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

	h.writeThumb(c, key, img)
}

func (h *Handler) writeThumb(c *gin.Context, key string, img image.Image) {
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
// @Param paths query string false "路径(可重复 path=a&path=b)"
// @Success 200 {file} binary
// @Router /api/filebrowser/download-tgz [get]
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
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fname))
	c.Header("Content-Type", "application/gzip")
	f, err := os.Open(tgzPath)
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "open failed", err)
		return
	}
	defer f.Close()
	io.Copy(c.Writer, f)
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

// Mkdir godoc
// @Summary 创建目录
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{path}"
// @Router /api/filebrowser/mkdir [post]
func (h *Handler) Mkdir(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	if mapErr(c, h.svc.Mkdir(req.Path)) {
		return
	}
	response.Ok(c, nil)
}

// Touch godoc
// @Summary 创建空文件（不覆盖）
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{path}"
// @Router /api/filebrowser/touch [post]
func (h *Handler) Touch(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	if mapErr(c, h.svc.Touch(req.Path)) {
		return
	}
	response.Ok(c, nil)
}

// Rename godoc
// @Summary 重命名文件/目录
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{path, new_name}"
// @Router /api/filebrowser/rename [post]
func (h *Handler) Rename(c *gin.Context) {
	var req struct {
		Path    string `json:"path"`
		NewName string `json:"new_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	if mapErr(c, h.svc.Rename(req.Path, req.NewName)) {
		return
	}
	response.Ok(c, nil)
}

// Move godoc
// @Summary 移动文件/目录到目标目录
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{paths, dest}"
// @Router /api/filebrowser/move [post]
func (h *Handler) Move(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths"`
		Dest  string   `json:"dest"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	count, err := h.svc.Move(req.Paths, req.Dest)
	if mapErr(c, err) {
		return
	}
	response.Ok(c, gin.H{"count": count})
}

// Copy godoc
// @Summary 复制文件/目录到目标目录
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{paths, dest}"
// @Router /api/filebrowser/copy [post]
func (h *Handler) Copy(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths"`
		Dest  string   `json:"dest"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	count, err := h.svc.Copy(req.Paths, req.Dest)
	if mapErr(c, err) {
		return
	}
	response.Ok(c, gin.H{"count": count})
}

// Delete godoc
// @Summary 删除文件/目录
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{paths}"
// @Router /api/filebrowser/delete [post]
func (h *Handler) Delete(c *gin.Context) {
	var req struct {
		Paths []string `json:"paths"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	count, err := h.svc.Delete(req.Paths)
	if mapErr(c, err) {
		return
	}
	response.Ok(c, gin.H{"count": count})
}

// ListTrash godoc
// @Summary 列出回收站
// @Tags 文件浏览
// @Router /api/filebrowser/trash [get]
func (h *Handler) ListTrash(c *gin.Context) {
	entries, err := h.svc.ListTrash()
	if mapErr(c, err) {
		return
	}
	response.Ok(c, entries)
}

// RestoreTrash godoc
// @Summary 从回收站恢复
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{trash_names}"
// @Router /api/filebrowser/trash/restore [post]
func (h *Handler) RestoreTrash(c *gin.Context) {
	var req struct {
		TrashNames []string `json:"trash_names"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	count, err := h.svc.RestoreTrash(req.TrashNames)
	if mapErr(c, err) {
		return
	}
	response.Ok(c, gin.H{"count": count})
}

// PermanentDelete godoc
// @Summary 从回收站永久删除
// @Tags 文件浏览
// @Accept json
// @Param body body object true "{trash_names}"
// @Router /api/filebrowser/trash/delete [post]
func (h *Handler) PermanentDelete(c *gin.Context) {
	var req struct {
		TrashNames []string `json:"trash_names"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, response.CodeInvalidParams, err.Error())
		return
	}
	count, err := h.svc.PermanentDelete(req.TrashNames)
	if mapErr(c, err) {
		return
	}
	response.Ok(c, gin.H{"count": count})
}

// Upload godoc
// @Summary 上传文件(不覆盖，支持多文件)
// @Description multipart/form-data 上传; 字段 file 为文件(可多个), path 为目标目录(相对 root, 空=root)
// @Tags 文件浏览
// @Accept multipart/form-data
// @Param path formData string false "目标目录(相对 root)"
// @Param file formData file true "文件(可多个)"
// @Success 200 {object} response.Response
// @Router /api/filebrowser/upload [post]
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
	overwrite := c.Query("overwrite") == "1" || c.PostForm("overwrite") == "1"

	results := make([]gin.H, 0, len(files))
	for _, fh := range files {
		r := h.uploadOne(dir, fh, overwrite)
		results = append(results, r)
	}
	response.Ok(c, results)
}

func (h *Handler) uploadOne(dir string, fh *multipart.FileHeader, overwrite bool) gin.H {
	if fh.Size > MaxUploadBytes {
		return gin.H{"name": fh.Filename, "error": "file too large"}
	}
	name := filepath.Base(fh.Filename)
	dst, finalRel, err := h.svc.CreateUploadFile(dir, name, overwrite)
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
