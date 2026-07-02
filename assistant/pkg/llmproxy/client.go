package llmproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"assistant/pkg/llm"
)

type ProxyClient struct {
	svc   *Service
	model string
}

func NewProxyClient(svc *Service) *ProxyClient {
	return &ProxyClient{
		svc:   svc,
		model: svc.Config().MiddleModel,
	}
}

func (c *ProxyClient) Provider() string { return "proxy" }

func (c *ProxyClient) Model() string { return c.model }

func (c *ProxyClient) Capabilities() llm.Capabilities {
	return llm.Capabilities{
		Supported: llm.CapabilityChat | llm.CapabilityStream | llm.CapabilityFunctionCall | llm.CapabilityVision,
	}
}

func (c *ProxyClient) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	payload := c.buildPayload(req, false)

	resp, err := c.svc.Forward(ctx, payload, req.Model)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, &llm.HTTPError{Code: resp.StatusCode, Message: string(b)}
	}

	var raw struct {
		Choices []struct {
			Message llm.Message `json:"message"`
		} `json:"choices"`
		Usage llm.TokenUsage `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(raw.Choices) == 0 {
		return nil, &llm.ProviderError{Message: "empty choices"}
	}

	return &llm.ChatResponse{
		Content:   raw.Choices[0].Message.Content,
		Role:      raw.Choices[0].Message.Role,
		ToolCalls: raw.Choices[0].Message.ToolCalls,
		Usage:     raw.Usage,
	}, nil
}

func (c *ProxyClient) StreamChat(ctx context.Context, req llm.ChatRequest, cb llm.StreamCallback) error {
	payload := c.buildPayload(req, true)

	resp, err := c.svc.Forward(ctx, payload, req.Model)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &llm.HTTPError{Code: resp.StatusCode, Message: string(b)}
	}

	err = llm.ReadSSE(ctx, resp.Body, func(data string) error {
		if data == "[DONE]" {
			return io.EOF
		}

		var raw struct {
			Choices []struct {
				Delta struct {
					Content   string         `json:"content"`
					Role      llm.Role       `json:"role"`
					ToolCalls []llm.ToolCall `json:"tool_calls"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Usage llm.TokenUsage `json:"usage"`
			Error any            `json:"error"`
		}

		if err = json.Unmarshal([]byte(data), &raw); err != nil {
			return err
		}
		if len(raw.Choices) == 0 {
			return nil
		}

		cb(llm.ChatResponse{
			Content:   raw.Choices[0].Delta.Content,
			Role:      raw.Choices[0].Delta.Role,
			ToolCalls: raw.Choices[0].Delta.ToolCalls,
			Usage:     raw.Usage,
		})

		if raw.Choices[0].FinishReason != nil && *raw.Choices[0].FinishReason != "" {
			return io.EOF
		}
		return nil
	})
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("proxy stream: %w", err)
	}
	return nil
}

func (c *ProxyClient) buildPayload(req llm.ChatRequest, stream bool) map[string]interface{} {
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
	if req.TopP > 0 {
		payload["top_p"] = req.TopP
	}
	if stream {
		payload["stream"] = true
	}
	if len(req.Tools) > 0 {
		payload["tools"] = req.Tools
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

func convertMessages(msgs []llm.Message) []map[string]interface{} {
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
