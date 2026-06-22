package kindle

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

func (m *Module) Name() string { return "kindle" }

func (m *Module) Register(r *gin.RouterGroup) {
	m.handler.Register(r)
}

func (m *Module) RegisterUI(r *gin.RouterGroup) {
	m.handler.RegisterUI(r)
}

func (m *Module) Middleware() []gin.HandlerFunc {
	return module.BaseMiddleware()
}
