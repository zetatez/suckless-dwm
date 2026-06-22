package filebrowser

import (
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

func (m *Module) Name() string { return "filebrowser" }

func (m *Module) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

func (m *Module) RegisterUI(r *gin.RouterGroup) {
	m.handler.RegisterUI(r)
}

func (m *Module) Middleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			if IsPublicPath(c.Query("path")) {
				c.Next()
				return
			}
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				c.Header("WWW-Authenticate", `Basic realm="filebrowser"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			cfg := psl.GetConfig()
			if username != cfg.Auth.Username || password != cfg.Auth.Password {
				c.Header("WWW-Authenticate", `Basic realm="filebrowser"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
