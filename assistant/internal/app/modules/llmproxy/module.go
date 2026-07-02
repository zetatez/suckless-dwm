package llmproxy

import (
	"net/http"

	"assistant/internal/app/module"
	"assistant/internal/bootstrap/psl"
	"assistant/pkg/llmproxy"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	svc     *llmproxy.Service
}

func NewModule() module.Module {
	svc := psl.GetProxyService()
	var handler *Handler
	if svc != nil && svc.HasProviders() {
		handler = NewHandler(svc)
	}
	return &Module{handler: handler, svc: svc}
}

func (m *Module) Name() string { return "llm" }

func (m *Module) Register(r *gin.RouterGroup) {
	if m.handler == nil {
		return
	}
	m.handler.Register(r)
}

func (m *Module) Middleware() []gin.HandlerFunc {
	if m.svc == nil {
		return module.BaseMiddleware()
	}
	token := m.svc.Config().AuthToken
	if token == "" {
		return module.BaseMiddleware()
	}
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			if c.GetHeader("Authorization") != "Bearer "+token {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
