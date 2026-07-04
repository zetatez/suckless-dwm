package llmproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"assistant/pkg/llm"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *llm.ProxyService
}

func NewHandler(svc *llm.ProxyService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	// OpenAI Chat Completions
	r.POST("/v1/chat/completions", h.chatCompletions)

	// OpenAI Responses API
	r.POST("/v1/responses", h.responsesHandler)
	r.POST("/responses", h.responsesHandlerLegacy)

	r.GET("/v1/models", h.ListModels)
	r.GET("/v1/models/:model", h.GetModel)
	r.GET("/status", h.Status)

	// Anthropic
	r.POST("/v1/messages", h.anthropicMessages)
}

// GetModel godoc
// @Summary 查询模型详情
// @Tags LLM代理
// @Param model path string true "模型ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/llmproxy/v1/models/{model} [get]
func (h *Handler) GetModel(c *gin.Context) {
	cfg := h.svc.Config()
	modelID := c.Param("model")

	if modelID == cfg.ProxiedModel {
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
// @Router /api/llmproxy/v1/models [get]
func (h *Handler) ListModels(c *gin.Context) {
	cfg := h.svc.Config()
	now := time.Now().Unix()
	seen := map[string]bool{cfg.ProxiedModel: true}
	data := []map[string]interface{}{
		{"id": cfg.ProxiedModel, "object": "model", "created": now, "owned_by": "assistant"},
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

// openaiResponse godoc
// @Summary OpenAI 格式聊天补全
// @Tags LLM代理
// @Success 200 {object} map[string]interface{}
// @Router /api/llmproxy/v1/chat/completions [post]
func (h *Handler) chatCompletions(c *gin.Context) {
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

	normalizeTools(reqMap)
	normalizeRoles(reqMap)

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

// responsesHandler handles OpenAI Responses API format (/v1/responses, /responses)
func (h *Handler) responsesHandler(c *gin.Context) {
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

	normalizeTools(reqMap)
	normalizeMessages(reqMap)
	normalizeRoles(reqMap)
	normalizeContent(reqMap)

	model, _ := reqMap["model"].(string)
	resp, err := h.svc.Forward(c.Request.Context(), reqMap, model)
	if err != nil {
		h.writeError(c, err)
		return
	}
	defer resp.Body.Close()

	stream, _ := reqMap["stream"].(bool)
	if stream {
		h.streamResponses(c, resp)
	} else {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json", oaiToResponses(raw))
	}
}

// responsesHandlerLegacy handles Codex CLI's /responses endpoint
func (h *Handler) responsesHandlerLegacy(c *gin.Context) {
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

	normalizeTools(reqMap)
	normalizeMessages(reqMap)
	normalizeRoles(reqMap)
	normalizeContent(reqMap)

	model, _ := reqMap["model"].(string)
	resp, err := h.svc.Forward(c.Request.Context(), reqMap, model)
	if err != nil {
		h.writeError(c, err)
		return
	}
	defer resp.Body.Close()

	stream, _ := reqMap["stream"].(bool)
	if stream {
		h.streamResponses(c, resp)
	} else {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json", oaiToResponses(raw))
	}
}

func oaiToResponses(body []byte) []byte {
	var oai struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Created int64  `json:"created"`
		Choices []struct {
			Message struct {
				Role             string `json:"role"`
				Content          string `json:"content"`
				ReasoningContent string `json:"reasoning_content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &oai); err != nil || len(oai.Choices) == 0 {
		return body
	}

	stopReason := oai.Choices[0].FinishReason
	switch stopReason {
	case "stop":
		stopReason = "end_turn"
	case "length":
		stopReason = "max_tokens"
	default:
		stopReason = "end_turn"
	}

	msg := oai.Choices[0].Message
	output := map[string]interface{}{
		"type": "message",
		"role": msg.Role,
		"content": []map[string]string{
			{"type": "output_text", "text": msg.Content},
		},
	}
	if msg.ReasoningContent != "" {
		output["reasoning_content"] = msg.ReasoningContent
	}

	resp := map[string]interface{}{
		"id":          "resp_" + oai.ID,
		"object":      "response",
		"created":     oai.Created,
		"model":       oai.Model,
		"output":      []map[string]interface{}{output},
		"stop_reason": stopReason,
	}
	if oai.Usage != nil {
		resp["usage"] = map[string]int{
			"input_tokens":  oai.Usage.PromptTokens,
			"output_tokens": oai.Usage.CompletionTokens,
		}
	}
	result, _ := json.Marshal(resp)
	return result
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

// streamResponses streams OpenAI Chat Completions SSE as Responses API SSE
func (h *Handler) streamResponses(c *gin.Context, resp *http.Response) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		raw, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json", oaiToResponses(raw))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	var itemID string
	var outputIndex int
	var contentBuf strings.Builder
	var reasoningBuf strings.Builder

	scanner := NewSSEScanner(resp.Body)
	for {
		_, data, err := scanner.Next()
		if err != nil {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Role             string `json:"role"`
					Content          string `json:"content"`
					ReasoningContent string `json:"reasoning_content"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil || len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta
		fr := chunk.Choices[0].FinishReason

		if (delta.Role == "assistant" || delta.ReasoningContent != "" || delta.Content != "") && itemID == "" {
			itemID = "item_" + fmt.Sprintf("%x", time.Now().UnixNano())
			added, _ := json.Marshal(map[string]interface{}{
				"type": "response.output_item.added",
				"item": map[string]interface{}{
					"id":      itemID,
					"type":    "message",
					"role":    "assistant",
					"content": []interface{}{},
				},
				"output_index": outputIndex,
			})
			fmt.Fprintf(c.Writer, "event: response.output_item.added\ndata: %s\n\n", added)
			flusher.Flush()
		}

		if delta.ReasoningContent != "" {
			reasoningBuf.WriteString(delta.ReasoningContent)
		}

		if delta.Content != "" {
			contentBuf.WriteString(delta.Content)
			textDelta, _ := json.Marshal(map[string]interface{}{
				"type":         "response.output_text.delta",
				"delta":        delta.Content,
				"item_id":      itemID,
				"output_index": outputIndex,
			})
			fmt.Fprintf(c.Writer, "event: response.output_text.delta\ndata: %s\n\n", textDelta)
			flusher.Flush()
		}

		if fr != "" {
			item := map[string]interface{}{
				"id":   itemID,
				"type": "message",
				"role": "assistant",
				"content": []map[string]string{
					{"type": "output_text", "text": contentBuf.String()},
				},
			}
			if reasoningBuf.Len() > 0 {
				item["reasoning_content"] = reasoningBuf.String()
			}
			done, _ := json.Marshal(map[string]interface{}{
				"type":         "response.output_item.done",
				"item":         item,
				"output_index": outputIndex,
			})
			fmt.Fprintf(c.Writer, "event: response.output_item.done\ndata: %s\n\n", done)

			completed, _ := json.Marshal(map[string]interface{}{
				"type": "response.completed",
				"response": map[string]interface{}{
					"id":     "resp_" + fmt.Sprintf("%x", time.Now().UnixNano()),
					"object": "response",
					"status": "completed",
					"output": []interface{}{},
				},
			})
			fmt.Fprintf(c.Writer, "event: response.completed\ndata: %s\n\n", completed)
			flusher.Flush()
			break
		}
	}
}

func (h *Handler) writeError(c *gin.Context, err error) {
	var httpErr *llm.HTTPError
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
// @Router /api/llmproxy/status [get]
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

// normalizeTools 将 tools 数组转换为 ark 兼容格式（只保留 function 类型）
func normalizeTools(req map[string]interface{}) {
	tools, ok := req["tools"].([]interface{})
	if !ok {
		return
	}
	var filtered []interface{}
	for _, t := range tools {
		tool, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		tt, _ := tool["type"].(string)
		if tt != "" && tt != "function" {
			continue
		}
		if _, hasFunc := tool["function"]; hasFunc {
			filtered = append(filtered, tool)
			continue
		}
		name, _ := tool["name"].(string)
		_, hasParams := tool["parameters"]
		if name == "" && !hasParams {
			continue
		}
		fn := map[string]interface{}{}
		if name != "" {
			fn["name"] = name
		}
		if desc, _ := tool["description"].(string); desc != "" {
			fn["description"] = desc
		}
		if hasParams {
			fn["parameters"] = tool["parameters"]
		}
		tool["type"] = "function"
		tool["function"] = fn
		delete(tool, "name")
		delete(tool, "description")
		delete(tool, "parameters")
		filtered = append(filtered, tool)
	}
	if len(filtered) == 0 {
		delete(req, "tools")
	} else {
		req["tools"] = filtered
	}
}

// normalizeMessages 将 Responses API 的 input 字段映射为 messages 字段
func normalizeMessages(req map[string]interface{}) {
	if _, ok := req["messages"]; ok {
		return
	}
	input, ok := req["input"]
	if !ok {
		return
	}
	switch v := input.(type) {
	case string:
		req["messages"] = []map[string]interface{}{
			{"role": "user", "content": v},
		}
	case []interface{}:
		req["messages"] = v
	}
	delete(req, "input")
}

// normalizeContent 将 Responses API 的消息内容块类型转为 Chat Completions 兼容格式
func normalizeContent(req map[string]interface{}) {
	msgs, ok := req["messages"].([]interface{})
	if !ok {
		return
	}
	for _, m := range msgs {
		msg, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		content, ok := msg["content"].([]interface{})
		if !ok {
			continue
		}
		for _, c := range content {
			block, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			t, _ := block["type"].(string)
			switch t {
			case "input_text":
				block["type"] = "text"
			case "input_image":
				block["type"] = "image_url"
			case "input_file":
				delete(block, "type")
			}
		}
	}
}

// normalizeRoles 将 Responses API 的角色名转换为 Chat Completions 兼容格式
func normalizeRoles(req map[string]interface{}) {
	msgs, ok := req["messages"].([]interface{})
	if !ok {
		return
	}
	for _, m := range msgs {
		msg, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		role, _ := msg["role"].(string)
		switch role {
		case "developer":
			msg["role"] = "system"
		}
	}
}
