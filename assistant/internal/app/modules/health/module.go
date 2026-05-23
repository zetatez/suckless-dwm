package health

import (
	"assistant/internal/app/module"

	"github.com/gin-gonic/gin"
)

type HealthModule struct {
	handler *HealthHandler
}

func NewHealthModule() module.Module {
	return &HealthModule{
		handler: NewHealthHandler(NewHealthService()),
	}
}

func (m *HealthModule) Name() string { return "health" }

func (m *HealthModule) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

func (m *HealthModule) Middleware() []gin.HandlerFunc {
	return module.BaseMiddleware()
}
