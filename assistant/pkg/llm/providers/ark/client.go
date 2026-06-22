package ark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"assistant/pkg/llm"
)

type imageURLContent struct {
	URL string `json:"url"`
}

type messageContent struct {
	Type     string           `json:"type"`
	Text     string           `json:"text,omitempty"`
	ImageURL *imageURLContent `json:"image_url,omitempty"`
}

type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *llm.BaseClient
}

func init() {
	llm.Register("ark", New)
}

func New(cfg llm.Config) (llm.Client, error) {
	if cfg.APIKey == "" {
		return nil, llm.ErrInvalidConfig
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://ark.cn-beijing.volces.com/api/coding/v3"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	model := cfg.Model
	if model == "" {
		model = "deepseek-v4-pro-260425"
	}

	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		client:  llm.NewBaseClient(baseURL, cfg),
	}, nil
}

func (c *Client) Provider() string { return "ark" }

func (c *Client) Model() string { return c.model }

func (c *Client) Capabilities() llm.Capabilities {
	return llm.Capabilities{Supported: llm.CapabilityChat | llm.CapabilityStream | llm.CapabilityFunctionCall | llm.CapabilityVision}
}

func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	messages := convertMessages(req.Messages)

	payload := map[string]any{
		"model":       c.getModel(req.Model),
		"messages":    messages,
		"temperature": req.Temperature,
	}

	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	if req.TopP > 0 {
		payload["top_p"] = req.TopP
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
	}

	resp, err := c.client.Do(ctx, "POST", "/chat/completions", payload, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 65536))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, &llm.HTTPError{Code: resp.StatusCode, Message: strings.TrimSpace(string(body))}
	}

	var raw struct {
		Choices []struct {
			Message llm.Message `json:"message"`
		} `json:"choices"`
		Usage llm.TokenUsage `json:"usage"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	if raw.Error != nil && raw.Error.Message != "" {
		return nil, &llm.ProviderError{Code: raw.Error.Code, Message: raw.Error.Message}
	}

	if len(raw.Choices) == 0 {
		return nil, &llm.ProviderError{Message: strings.TrimSpace(string(body))}
	}

	return &llm.ChatResponse{
		Content:   raw.Choices[0].Message.Content,
		Role:      raw.Choices[0].Message.Role,
		ToolCalls: raw.Choices[0].Message.ToolCalls,
		Usage:     raw.Usage,
	}, nil
}

func (c *Client) StreamChat(ctx context.Context, req llm.ChatRequest, cb llm.StreamCallback) error {
	messages := convertMessages(req.Messages)

	payload := map[string]any{
		"model":       c.getModel(req.Model),
		"messages":    messages,
		"temperature": req.Temperature,
		"stream":      true,
	}

	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	if req.TopP > 0 {
		payload["top_p"] = req.TopP
	}
	if len(req.Tools) > 0 {
		payload["tools"] = req.Tools
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.HTTPClient().Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &llm.HTTPError{Code: resp.StatusCode, Message: strings.TrimSpace(string(b))}
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

		if err := json.Unmarshal([]byte(data), &raw); err != nil {
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
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("ark stream: %w", err)
	}
	return nil
}

func (c *Client) getModel(model string) string {
	if model != "" {
		return model
	}
	return c.model
}

func convertMessages(msgs []llm.Message) []map[string]any {
	result := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		if m.ImageBase64 != "" {
			content := []messageContent{
				{Type: "text", Text: m.Content},
				{Type: "image_url", ImageURL: &imageURLContent{URL: "data:image/jpeg;base64," + m.ImageBase64}},
			}
			result = append(result, map[string]any{
				"role":    m.Role,
				"content": content,
			})
		} else {
			result = append(result, map[string]any{
				"role":    m.Role,
				"content": m.Content,
			})
		}
	}
	return result
}
