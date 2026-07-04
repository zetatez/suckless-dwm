package svc

import (
	"assistant/internal/app/module"
	"assistant/internal/bootstrap/psl"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule() module.Module {
	return &Module{
		handler: NewHandler(NewService()),
	}
}

func (m *Module) Name() string { return "svr" }

func (m *Module) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

func (m *Module) Middleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				c.Header("WWW-Authenticate", `Basic realm="cmd"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			cfg := psl.GetConfig()
			if username != cfg.App.RootUsername || password != cfg.App.RootPassword {
				c.Header("WWW-Authenticate", `Basic realm="cmd"`)
				c.AbortWithStatusJSON(401, gin.H{"code": 40101, "message": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
