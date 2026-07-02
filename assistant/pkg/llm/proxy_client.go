package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type ProxyClient struct {
	svc *ProxyService
}

func NewProxyClient(svc *ProxyService) *ProxyClient {
	return &ProxyClient{svc: svc}
}

func (c *ProxyClient) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	payload := c.buildPayload(req)

	resp, err := c.svc.Forward(ctx, payload, req.Model)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, &HTTPError{Code: resp.StatusCode, Message: string(b)}
	}

	var raw struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage TokenUsage `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(raw.Choices) == 0 {
		return nil, fmt.Errorf("empty choices")
	}

	return &ChatResponse{
		Content:   raw.Choices[0].Message.Content,
		Role:      raw.Choices[0].Message.Role,
		ToolCalls: raw.Choices[0].Message.ToolCalls,
		Usage:     raw.Usage,
	}, nil
}

func (c *ProxyClient) buildPayload(req ChatRequest) map[string]interface{} {
	cfg := c.svc.Config()
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    convertMessages(req.Messages),
		"temperature": req.Temperature,
	}
	if cfg.Temperature > 0 && req.Temperature == 0 {
		payload["temperature"] = cfg.Temperature
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	return payload
}

type imageURLContent struct {
	URL string `json:"url"`
}

type messageContent struct {
	Type     string           `json:"type"`
	Text     string           `json:"text,omitempty"`
	ImageURL *imageURLContent `json:"image_url,omitempty"`
}

func convertMessages(msgs []Message) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(msgs))
	for _, m := range msgs {
		if m.ImageBase64 != "" {
			content := []messageContent{
				{Type: "text", Text: m.Content},
				{Type: "image_url", ImageURL: &imageURLContent{URL: "data:image/jpeg;base64," + m.ImageBase64}},
			}
			result = append(result, map[string]interface{}{
				"role":    m.Role,
				"content": content,
			})
		} else {
			result = append(result, map[string]interface{}{
				"role":    m.Role,
				"content": m.Content,
			})
		}
	}
	return result
}
