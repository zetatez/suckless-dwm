package llmproxy

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"assistant/pkg/llmproxy"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *llmproxy.Service
}

func NewHandler(svc *llmproxy.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	r.POST("/v1/chat/completions", h.ChatCompletions)
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

func copyHeaders(c *gin.Context, h http.Header) {
	for k, v := range h {
		for _, val := range v {
			c.Writer.Header().Add(k, val)
		}
	}
}

// ChatCompletions godoc
// @Summary 聊天补全(自动路由到可用供应商)
// @Tags LLM代理
// @Param request body object true "请求体"
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/chat/completions [post]
func (h *Handler) ChatCompletions(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}

	var reqMap map[string]interface{}
	if err := json.Unmarshal(rawBody, &reqMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	requestedModel, _ := reqMap["model"].(string)
	reqStream, _ := reqMap["stream"].(bool)

	resp, err := h.svc.Forward(c.Request.Context(), reqMap, requestedModel)
	if err != nil {
		if err == llmproxy.ErrNoProvider {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no available provider"})
			return
		}
		if err == llmproxy.ErrAllProvidersFailed {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "all providers failed"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		copyHeaders(c, resp.Header)
		raw, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), raw)
		return
	}

	copyHeaders(c, resp.Header)

	if !reqStream {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, resp.Header.Get("Content-Type"), raw)
		return
	}

	c.Status(http.StatusOK)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, resp.Header.Get("Content-Type"), raw)
		return
	}

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
