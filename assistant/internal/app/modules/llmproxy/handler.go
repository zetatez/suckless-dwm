package llmproxy

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	llmproxysvc "assistant/pkg/llmproxy"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/v1/chat/completions", h.ChatCompletions)
	r.POST("/v1/messages", h.AnthropicMessages)
	r.GET("/v1/models", h.ListModels)
	r.GET("/v1/models/:model", h.GetModel)
	r.GET("/status", h.Status)
}

// GetModel godoc
// @Summary 查询模型详情
// @Tags LLM代理
// @Param model path string true "模型ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/models/{model} [get]
func (h *Handler) GetModel(c *gin.Context) {
	cfg := h.svc.Config()
	modelID := c.Param("model")

	if modelID == cfg.MiddleModel {
		c.JSON(http.StatusOK, map[string]interface{}{
			"id": modelID, "object": "model", "created": time.Now().Unix(), "owned_by": "assistant",
		})
		return
	}
	for _, p := range cfg.Providers {
		if hasModel(p.Models, modelID) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"id": modelID, "object": "model", "created": time.Now().Unix(), "owned_by": p.Name,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "model not found", "type": "invalid_request_error", "code": "model_not_found"}})
}

// ListModels godoc
// @Summary 列出可用模型
// @Tags LLM代理
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/models [get]
func (h *Handler) ListModels(c *gin.Context) {
	cfg := h.svc.Config()
	now := time.Now().Unix()
	seen := map[string]bool{cfg.MiddleModel: true}
	data := []map[string]interface{}{
		{"id": cfg.MiddleModel, "object": "model", "created": now, "owned_by": "assistant"},
	}
	for _, p := range cfg.Providers {
		for _, m := range p.Models {
			if seen[m] {
				continue
			}
			seen[m] = true
			data = append(data, map[string]interface{}{
				"id": m, "object": "model", "created": now, "owned_by": p.Name,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"object": "list", "data": data})
}

// ChatCompletions godoc
// @Summary 聊天补全(自动路由到可用供应商)
// @Tags LLM代理
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/chat/completions [post]
func (h *Handler) ChatCompletions(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}

	var reqMap map[string]interface{}
	if err := json.Unmarshal(body, &reqMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	model, _ := reqMap["model"].(string)
	resp, err := h.svc.Forward(c.Request.Context(), reqMap, model)
	if err != nil {
		h.writeError(c, err)
		return
	}
	defer resp.Body.Close()

	stream, _ := reqMap["stream"].(bool)
	if stream {
		h.streamOpenAI(c, resp)
	} else {
		raw, _ := io.ReadAll(resp.Body)
		for k, v := range resp.Header {
			c.Header(k, v[0])
		}
		c.Data(http.StatusOK, resp.Header.Get("Content-Type"), raw)
	}
}

func (h *Handler) streamOpenAI(c *gin.Context, resp *http.Response) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, resp.Header.Get("Content-Type"), raw)
		return
	}

	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	c.Status(http.StatusOK)

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			c.Writer.Write(buf[:n])
			flusher.Flush()
		}
		if err != nil {
			break
		}
	}
}

func (h *Handler) writeError(c *gin.Context, err error) {
	var httpErr *llmproxysvc.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
}

// Status godoc
// @Summary 查看供应商状态
// @Tags LLM代理
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/status [get]
func (h *Handler) Status(c *gin.Context) {
	active := h.svc.ActiveProvider()
	resp := map[string]interface{}{
		"providers": h.svc.ProviderStatuses(),
	}
	if active != nil {
		resp["active_provider"] = active.Name
	}
	c.JSON(http.StatusOK, resp)
}

func hasModel(models []string, target string) bool {
	for _, m := range models {
		if m == target {
			return true
		}
	}
	return false
}
