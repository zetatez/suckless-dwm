package llmproxy

import (
	"net/http"

	"assistant/internal/app/module"
	"assistant/internal/bootstrap/psl"
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
	router  *Router
}

func NewModule() module.Module {
	cfg := psl.GetConfig().LLMProxy
	hasProviders := len(cfg.Providers) > 0
	var router *Router
	var handler *Handler
	if hasProviders {
		router = NewRouter()
		handler = NewHandler(router)
	}
	return &Module{
		handler: handler,
		router:  router,
	}
}

func (m *Module) Name() string { return "llm" }

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
