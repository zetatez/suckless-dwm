package health

import (
	"assistant/internal/app/module"

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

func (m *Module) Name() string { return "health" }

func (m *Module) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

func (m *Module) RegisterUI(r *gin.RouterGroup) {} // no UI

func (m *Module) Middleware() []gin.HandlerFunc {
	return module.BaseMiddleware()
}
