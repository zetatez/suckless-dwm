package llmproxy

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	llmproxysvc "assistant/pkg/llmproxy"

	"github.com/gin-gonic/gin"
)

type anthropicReq struct {
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
	Messages  []anthropicMsg `json:"messages"`
	System    string         `json:"system,omitempty"`
	Stream    bool           `json:"stream,omitempty"`
}

type anthropicMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func anthropicToOpenAI(body []byte) ([]byte, error) {
	var areq anthropicReq
	if err := json.Unmarshal(body, &areq); err != nil {
		return nil, err
	}
	omap := map[string]interface{}{
		"model":      areq.Model,
		"max_tokens": areq.MaxTokens,
		"stream":     areq.Stream,
	}
	var omsgs []map[string]interface{}
	if areq.System != "" {
		omsgs = append(omsgs, map[string]interface{}{"role": "system", "content": areq.System})
	}
	for _, m := range areq.Messages {
		if m.Role == "" || m.Content == "" {
			continue
		}
		omsgs = append(omsgs, map[string]interface{}{"role": m.Role, "content": m.Content})
	}
	omap["messages"] = omsgs
	return json.Marshal(omap)
}

func openAIToAnthropic(body []byte) ([]byte, error) {
	var ores struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &ores); err != nil {
		return nil, err
	}

	ares := map[string]interface{}{
		"id":            "msg_" + fmt.Sprintf("%x", time.Now().UnixNano()),
		"type":          "message",
		"role":          "assistant",
		"model":         ores.Model,
		"content":       []map[string]string{{"type": "text", "text": ""}},
		"stop_reason":   "end_turn",
		"stop_sequence": nil,
	}
	if len(ores.Choices) > 0 {
		if ores.Choices[0].Message.Content != "" {
			ares["content"] = []map[string]string{
				{"type": "text", "text": ores.Choices[0].Message.Content},
			}
		}
		switch ores.Choices[0].FinishReason {
		case "stop":
			ares["stop_reason"] = "end_turn"
		case "length":
			ares["stop_reason"] = "max_tokens"
		}
	}
	if ores.Usage != nil {
		ares["usage"] = map[string]int{
			"input_tokens":  ores.Usage.PromptTokens,
			"output_tokens": ores.Usage.CompletionTokens,
		}
	}
	return json.Marshal(ares)
}

// AnthropicMessages godoc
// @Summary Anthropic 格式聊天补全
// @Tags LLM代理
// @Param request body anthropicReq true "Anthropic 格式请求体"
// @Success 200 {object} map[string]interface{}
// @Router /api/llm/v1/messages [post]
func (h *Handler) AnthropicMessages(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeAnthropicError(c, http.StatusBadRequest, "invalid_request_error", "read body failed")
		return
	}

	var areq anthropicReq
	if err := json.Unmarshal(body, &areq); err != nil {
		writeAnthropicError(c, http.StatusBadRequest, "invalid_request_error", "invalid json")
		return
	}
	if len(areq.Messages) == 0 {
		writeAnthropicError(c, http.StatusBadRequest, "invalid_request_error", "messages is required")
		return
	}

	oaiBody, err := anthropicToOpenAI(body)
	if err != nil {
		writeAnthropicError(c, http.StatusInternalServerError, "internal_error", "convert failed")
		return
	}

	var reqMap map[string]interface{}
	json.Unmarshal(oaiBody, &reqMap)

	resp, err := h.svc.Forward(c.Request.Context(), reqMap, areq.Model)
	if err != nil {
		var httpErr *llmproxysvc.HTTPError
		if errors.As(err, &httpErr) {
			writeAnthropicError(c, httpErr.Code, "api_error", httpErr.Message)
		} else {
			writeAnthropicError(c, http.StatusServiceUnavailable, "overloaded_error", err.Error())
		}
		return
	}
	defer resp.Body.Close()

	if areq.Stream {
		h.anthropicStream(c, resp)
	} else {
		raw, _ := io.ReadAll(resp.Body)
		ares, err := openAIToAnthropic(raw)
		if err != nil {
			writeAnthropicError(c, http.StatusInternalServerError, "internal_error", "convert response failed")
			return
		}
		c.Data(http.StatusOK, "application/json", ares)
	}
}

func (h *Handler) anthropicStream(c *gin.Context, resp *http.Response) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		raw, _ := io.ReadAll(resp.Body)
		ares, _ := openAIToAnthropic(raw)
		c.Data(http.StatusOK, "application/json", ares)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	msgID := "msg_" + fmt.Sprintf("%x", time.Now().UnixNano())
	var contentSent bool

	type oaiChunk struct {
		Choices []struct {
			Delta struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"delta"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	// message_start
	startEvent, _ := json.Marshal(map[string]interface{}{
		"type": "message_start",
		"message": map[string]interface{}{
			"id":      msgID,
			"type":    "message",
			"role":    "assistant",
			"content": []interface{}{},
		},
	})
	fmt.Fprintf(c.Writer, "event: message_start\ndata: %s\n\n", startEvent)
	flusher.Flush()

	scanner := NewSSEScanner(resp.Body)
	for {
		_, data, err := scanner.Next()
		if err != nil {
			break
		}

		var chunk oaiChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil || len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta

		if delta.Role == "assistant" && !contentSent {
			cbStart, _ := json.Marshal(map[string]interface{}{
				"type":  "content_block_start",
				"index": 0,
				"content_block": map[string]interface{}{
					"type": "text",
					"text": "",
				},
			})
			fmt.Fprintf(c.Writer, "event: content_block_start\ndata: %s\n\n", cbStart)
			contentSent = true
			flusher.Flush()
		}

		if delta.Content != "" {
			cbDelta, _ := json.Marshal(map[string]interface{}{
				"type":  "content_block_delta",
				"index": 0,
				"delta": map[string]string{"type": "text", "text": delta.Content},
			})
			fmt.Fprintf(c.Writer, "event: content_block_delta\ndata: %s\n\n", cbDelta)
			flusher.Flush()
		}

		if chunk.Choices[0].FinishReason != "" {
			fr := chunk.Choices[0].FinishReason
			stopReason := "end_turn"
			if fr == "length" {
				stopReason = "max_tokens"
			}

			cbStop, _ := json.Marshal(map[string]interface{}{
				"type": "content_block_stop", "index": 0,
			})
			fmt.Fprintf(c.Writer, "event: content_block_stop\ndata: %s\n\n", cbStop)

			usage := map[string]int{}
			if chunk.Usage != nil {
				usage["input_tokens"] = chunk.Usage.PromptTokens
				usage["output_tokens"] = chunk.Usage.CompletionTokens
			}
			msgDelta, _ := json.Marshal(map[string]interface{}{
				"type": "message_delta",
				"delta": map[string]interface{}{
					"stop_reason":   stopReason,
					"stop_sequence": nil,
				},
				"usage": usage,
			})
			fmt.Fprintf(c.Writer, "event: message_delta\ndata: %s\n\n", msgDelta)
			fmt.Fprintf(c.Writer, "event: message_stop\ndata: {}\n\n")
			flusher.Flush()
			break
		}
	}
}

type sseScanner struct {
	reader  *bufio.Reader
	dataBuf strings.Builder
}

func NewSSEScanner(r io.Reader) *sseScanner {
	return &sseScanner{reader: bufio.NewReader(r)}
}

func (s *sseScanner) Next() (event, data string, err error) {
	s.dataBuf.Reset()
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			return s.dataBuf.String(), s.dataBuf.String(), nil
		}
		if strings.HasPrefix(line, "data: ") {
			text := strings.TrimPrefix(line, "data: ")
			if s.dataBuf.Len() > 0 {
				s.dataBuf.WriteString("\n")
			}
			s.dataBuf.WriteString(text)
		}
	}
}

func writeAnthropicError(c *gin.Context, code int, etype, msg string) {
	c.JSON(code, map[string]interface{}{
		"error": map[string]interface{}{
			"type":    etype,
			"message": msg,
		},
	})
}
