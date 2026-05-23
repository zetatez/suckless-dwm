package module

import "github.com/gin-gonic/gin"

type Module interface {
	Name() string
	Register(r *gin.RouterGroup)
	Middleware() []gin.HandlerFunc
}

func BaseMiddleware() []gin.HandlerFunc {
	return nil
}
