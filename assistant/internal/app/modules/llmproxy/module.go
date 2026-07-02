package llmproxy

import (
	"net/http"

	"assistant/internal/app/module"
	"assistant/internal/bootstrap/psl"
	"assistant/pkg/llm"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule() module.Module {
	cfg := psl.GetConfig().LLMProxy
	if len(cfg.Providers) == 0 {
		return &Module{}
	}
	return &Module{
		handler: NewHandler(llm.NewProxyService(cfg)),
	}
}

func (m *Module) Name() string { return "llmproxy" }

func (m *Module) Register(r *gin.RouterGroup) {
	if m.handler == nil {
		return
	}
	m.handler.Register(r)
}

func (m *Module) RegisterUI(r *gin.RouterGroup) {}

func (m *Module) Middleware() []gin.HandlerFunc {
	cfg := psl.GetConfig().LLMProxy
	if cfg.AuthToken == "" {
		return module.BaseMiddleware()
	}
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			if c.GetHeader("Authorization") != "Bearer "+cfg.AuthToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			c.Next()
		},
	}
}
