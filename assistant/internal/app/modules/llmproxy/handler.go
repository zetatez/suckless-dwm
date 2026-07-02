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
	svc *Service
}

func NewHandler(svc *Service) *Handler {
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
	cfg := psl.GetConfig().LLMProxy
	modelID := c.Param("model")

	if modelID == cfg.MiddleModel {
		c.JSON(http.StatusOK, map[string]interface{}{
			"id": modelID, "object": "model", "created": time.Now().Unix(), "owned_by": "assistant",
		})
		return
	}
	for _, p := range h.svc.providers {
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
	for _, p := range h.svc.providers {
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

	cfg := psl.GetConfig().LLMProxy
	var p *providerState
	var modelSpecific bool
	if requestedModel == cfg.MiddleModel {
		p = h.svc.selectProvider()
	} else {
		p = h.svc.selectProviderByModel(requestedModel)
		if p == nil {
			p = h.svc.selectProvider()
		} else {
			modelSpecific = true
		}
	}

	if p == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no available provider"})
		return
	}

	reqMap["model"] = resolveModel(p, requestedModel)
	body, _ := json.Marshal(reqMap)

	if reqStream {
		h.streamResponse(c, p, body, modelSpecific, requestedModel, reqMap)
	} else {
		h.syncResponse(c, p, body, modelSpecific, requestedModel, reqMap)
	}
}

func resolveModel(p *providerState, requested string) string {
	if hasModel(p.Models, requested) {
		return requested
	}
	return p.Models[0]
}

func (h *Handler) syncResponse(c *gin.Context, p *providerState, body []byte, modelSpecific bool, originalModel string, reqMap map[string]interface{}) {
	for {
		resp, err := h.svc.ForwardChat(c.Request.Context(), p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			logError("forward to %s failed: %v", p.Name, err)
			h.svc.MarkOffline(p)
			p = h.findNext(p, modelSpecific, originalModel)
			if p == nil {
				break
			}
			reqMap["model"] = resolveModel(p, originalModel)
			body, _ = json.Marshal(reqMap)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			copyHeaders(c, resp.Header)
			raw, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), raw)
			h.svc.NotifyActivity()
			return
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if p.PlanType == "fixed" {
			h.svc.MarkExhausted(p)
		} else {
			h.svc.MarkOffline(p)
		}

		pNext := h.findNext(p, modelSpecific, originalModel)
		if pNext == nil {
			c.JSON(resp.StatusCode, gin.H{"error": "provider failed", "provider": p.Name})
			return
		}
		p = pNext
		reqMap["model"] = resolveModel(p, originalModel)
		body, _ = json.Marshal(reqMap)
	}

	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "all providers failed"})
}

func (h *Handler) streamResponse(c *gin.Context, p *providerState, body []byte, modelSpecific bool, originalModel string, reqMap map[string]interface{}) {
	for {
		resp, err := h.svc.ForwardChat(c.Request.Context(), p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			logError("stream forward to %s failed: %v", p.Name, err)
			h.svc.MarkOffline(p)
			next := h.findNext(p, modelSpecific, originalModel)
			if next == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "all providers failed"})
				return
			}
			p = next
			reqMap["model"] = resolveModel(p, originalModel)
			body, _ = json.Marshal(reqMap)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if p.PlanType == "fixed" {
				h.svc.MarkExhausted(p)
			} else {
				h.svc.MarkOffline(p)
			}
			next := h.findNext(p, modelSpecific, originalModel)
			if next == nil {
				c.JSON(resp.StatusCode, gin.H{"error": "provider failed", "provider": p.Name})
				return
			}
			p = next
			reqMap["model"] = resolveModel(p, originalModel)
			body, _ = json.Marshal(reqMap)
			continue
		}

		copyHeaders(c, resp.Header)
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
		h.svc.NotifyActivity()

		if c.Request.Context().Err() != nil {
			logError("stream %s cancelled by client", p.Name)
		}
		return
	}
}

func (h *Handler) findNext(current *providerState, modelSpecific bool, originalModel string) *providerState {
	h.svc.mu.RLock()
	defer h.svc.mu.RUnlock()
	for _, p := range h.svc.providers {
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
