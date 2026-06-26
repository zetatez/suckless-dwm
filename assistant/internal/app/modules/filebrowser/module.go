package filebrowser

import (
	"strings"

	"assistant/internal/app/module"
	"assistant/internal/bootstrap/psl"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule() module.Module {
	return &Module{handler: NewHandler(NewService())}
}

func (m *Module) Name() string { return "files" }

func (m *Module) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

// isPublicEndpoint UI 入口本身始终免鉴权(只是静态 HTML)。
func isPublicEndpoint(c *gin.Context) bool {
	p := c.Request.URL.Path
	return p == "/api/files" || p == "/api/files/" || strings.HasSuffix(p, "/files/ui")
}

// extractPath 仅从 query 取 path(避免提前消费 multipart body)。
// 对于上传到 public 目录，前端需要把 path 放到 URL query 中。
func extractPath(c *gin.Context) string {
	return c.Query("path")
}

func (m *Module) Middleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			if isPublicEndpoint(c) {
				c.Next()
				return
			}
			if IsPublicPath(extractPath(c)) {
				c.Next()
				return
			}
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				c.Header("WWW-Authenticate", `Basic realm="files"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			cfg := psl.GetConfig()
			if username != cfg.Auth.Username || password != cfg.Auth.Password {
				c.Header("WWW-Authenticate", `Basic realm="files"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
