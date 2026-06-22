package llmproxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"assistant/internal/bootstrap/psl"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	router *Router
}

func NewHandler(router *Router) *Handler {
	return &Handler{router: router}
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
	cfg := psl.GetConfig().LLMProxy
	modelID := c.Param("model")

	if modelID == cfg.MiddleModel {
		c.JSON(http.StatusOK, map[string]interface{}{
			"id": modelID, "object": "model", "created": time.Now().Unix(), "owned_by": "assistant",
		})
		return
	}
	for _, p := range h.router.providers {
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
	cfg := psl.GetConfig().LLMProxy
	now := time.Now().Unix()
	seen := map[string]bool{cfg.MiddleModel: true}
	data := []map[string]interface{}{
		{"id": cfg.MiddleModel, "object": "model", "created": now, "owned_by": "assistant"},
	}
	for _, p := range h.router.providers {
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

type proxyChatRequest struct {
	Model    string `json:"model"`
	Messages []any  `json:"messages"`
	Stream   bool   `json:"stream,omitempty"`
}

// ChatCompletions godoc
// @Summary 聊天补全(自动路由到可用供应商)
// @Tags LLM代理
// @Param request body proxyChatRequest true "请求体"
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/chat/completions [post]
func (h *Handler) ChatCompletions(c *gin.Context) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}
	var req proxyChatRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	cfg := psl.GetConfig().LLMProxy
	requestedModel := req.Model

	var p *providerState
	var modelSpecific bool
	if requestedModel == cfg.MiddleModel {
		p = h.router.selectProvider()
	} else {
		p = h.router.selectProviderByModel(requestedModel)
		if p == nil {
			p = h.router.selectProvider()
		} else {
			modelSpecific = true
		}
	}

	if p == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no available provider"})
		return
	}

	body := patchModel(rawBody, resolveModel(p, requestedModel))

	if req.Stream {
		h.streamResponse(c, p, body, modelSpecific, requestedModel)
	} else {
		h.syncResponse(c, p, body, modelSpecific, requestedModel)
	}
}

func resolveModel(p *providerState, requested string) string {
	if hasModel(p.Models, requested) {
		return requested
	}
	return p.Models[0]
}

func (h *Handler) syncResponse(c *gin.Context, p *providerState, body []byte, modelSpecific bool, originalModel string) {
	for {
		resp, err := h.router.ForwardChat(c.Request.Context(), p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			logError("forward to %s failed: %v", p.Name, err)
			h.router.MarkOffline(p.Name)
			p = h.findNext(p, modelSpecific, originalModel)
			if p == nil {
				break
			}
			body = patchModel(body, resolveModel(p, originalModel))
			continue
		}

		if resp.StatusCode == http.StatusOK {
			for k, v := range resp.Header {
				c.Header(k, v[0])
			}
			raw, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), raw)
			h.router.NotifyActivity()
			return
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if p.PlanType == "fixed" {
			h.router.MarkExhausted(p.Name)
		} else {
			h.router.MarkOffline(p.Name)
		}

		pNext := h.findNext(p, modelSpecific, originalModel)
		if pNext == nil {
			c.JSON(resp.StatusCode, gin.H{"error": "provider failed", "provider": p.Name})
			return
		}
		p = pNext
		body = patchModel(body, resolveModel(p, originalModel))
	}

	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "all providers failed"})
}

func (h *Handler) streamResponse(c *gin.Context, p *providerState, body []byte, modelSpecific bool, originalModel string) {
	resp, err := h.router.ForwardChat(c.Request.Context(), p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		logError("stream forward to %s failed: %v", p.Name, err)
		h.router.MarkOffline(p.Name)
		if next := h.findNext(p, modelSpecific, originalModel); next != nil {
			body = patchModel(body, resolveModel(next, originalModel))
			h.streamResponse(c, next, body, modelSpecific, originalModel)
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "all providers failed"})
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		h.router.MarkExhausted(p.Name)
		if next := h.findNext(p, modelSpecific, originalModel); next != nil {
			body = patchModel(body, resolveModel(next, originalModel))
			h.streamResponse(c, next, body, modelSpecific, originalModel)
		} else {
			c.JSON(resp.StatusCode, gin.H{"error": "provider failed", "provider": p.Name})
		}
		return
	}

	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
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
	resp.Body.Close()
	h.router.NotifyActivity()

	if c.Request.Context().Err() != nil {
		logError("stream %s cancelled by client", p.Name)
	}
}

func (h *Handler) findNext(current *providerState, modelSpecific bool, originalModel string) *providerState {
	h.router.mu.RLock()
	defer h.router.mu.RUnlock()
	for _, p := range h.router.providers {
		if p == current {
			continue
		}
		if p.status != StatusAvailable {
			continue
		}
		if modelSpecific && !hasModel(p.Models, originalModel) {
			continue
		}
		return p
	}
	return nil
}

func patchModel(body []byte, model string) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return body
	}
	m["model"] = model
	patched, _ := json.Marshal(m)
	return patched
}

func (h *Handler) Status(c *gin.Context) {
	active := h.router.ActiveProvider()
	resp := map[string]interface{}{
		"providers": h.router.ProviderStatuses(),
	}
	if active != nil {
		resp["active_provider"] = active.Name
	}
	c.JSON(http.StatusOK, resp)
}
