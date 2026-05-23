package health

import (
	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	svc *HealthService
}

func NewHealthHandler(svc *HealthService) *HealthHandler {
	return &HealthHandler{svc: svc}
}

func (h *HealthHandler) Register(r *gin.RouterGroup) {
	r.GET("", h.Health)
}

// Health godoc
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags 健康检查
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Failure 500 {object} response.Response "服务不健康"
// @Router /api/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	data, err := h.svc.Health()
	if err != nil {
		response.ErrWithInternal(c, response.CodeServerError, "health check failed", err)
		return
	}
	if data.Status != "ok" {
		response.Err(c, response.CodeServerError, data.Status)
		return
	}
	response.Ok(c, data)
}
