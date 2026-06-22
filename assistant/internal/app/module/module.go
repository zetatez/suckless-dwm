package module

import "github.com/gin-gonic/gin"

type Module interface {
	Name() string
	Register(r *gin.RouterGroup)
	RegisterUI(r *gin.RouterGroup)
	Middleware() []gin.HandlerFunc
}

func BaseMiddleware() []gin.HandlerFunc {
	return nil
}

// NoopUI 默认 RegisterUI 实现，适合纯 API 模块。
func NoopUI(r *gin.RouterGroup) {} // nop
