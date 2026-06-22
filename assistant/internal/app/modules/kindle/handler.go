package kindle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {}

func (h *Handler) RegisterUI(r *gin.RouterGroup) {
	r.GET("", func(c *gin.Context) {
		if err := renderUI(h.svc, c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})
	r.GET("/content", func(c *gin.Context) {
		if err := renderContent(h.svc, c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})
}
