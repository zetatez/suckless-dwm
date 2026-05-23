package minimax

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

type minimaxContentBlock struct {
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *llm.BaseClient
}

func init() {
	llm.Register("minimax", New)
}

func New(cfg llm.Config) (llm.Client, error) {
	if cfg.APIKey == "" {
		return nil, llm.ErrInvalidConfig
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.minimax.chat/v1"
	}

	model := cfg.Model
	if model == "" {
		model = "abab6.5s-chat"
	}

	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		client:  llm.NewBaseClient(baseURL, cfg),
	}, nil
}

func (c *Client) Provider() string { return "minimax" }

func (c *Client) Model() string { return c.model }

func (c *Client) Capabilities() llm.Capabilities {
	return llm.Capabilities{Supported: llm.CapabilityChat | llm.CapabilityStream}
}

func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	messages := convertMessages(req.Messages)

	payload := map[string]any{
		"model":              c.getModel(req.Model),
		"messages":           messages,
		"temperature":        req.Temperature,
		"tokens_to_generate": 4096,
	}

	if req.MaxTokens > 0 {
		payload["tokens_to_generate"] = req.MaxTokens
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
	}

	resp, err := c.client.Do(ctx, "POST", "/v1/text/chatcompletion_v2", payload, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
				Role    string `json:"role"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			TotalTokens      int `json:"total_tokens"`
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
		Error struct {
			Code    int    `json:"status_code"`
			Message string `json:"status_msg"`
		} `json:"base_resp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if raw.Error.Code != 0 {
		return nil, &llm.ProviderError{Code: fmt.Sprintf("%d", raw.Error.Code), Message: raw.Error.Message}
	}

	if len(raw.Choices) == 0 {
		return nil, llm.ErrMaxRetries
	}

	return &llm.ChatResponse{
		Content: raw.Choices[0].Message.Content,
		Role:    llm.Role(raw.Choices[0].Message.Role),
		Usage: llm.TokenUsage{
			TotalTokens:      raw.Usage.TotalTokens,
			PromptTokens:     raw.Usage.PromptTokens,
			CompletionTokens: raw.Usage.CompletionTokens,
		},
	}, nil
}

func (c *Client) StreamChat(ctx context.Context, req llm.ChatRequest, cb llm.StreamCallback) error {
	payload := map[string]any{
		"model":              c.getModel(req.Model),
		"messages":           req.Messages,
		"temperature":        req.Temperature,
		"tokens_to_generate": 4096,
		"stream":             true,
	}
	if req.MaxTokens > 0 {
		payload["tokens_to_generate"] = req.MaxTokens
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
		"Accept":        "text/event-stream",
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/text/chatcompletion_v2", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

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
					Content string `json:"content"`
					Role    string `json:"role"`
				} `json:"delta"`
				Message struct {
					Content string `json:"content"`
					Role    string `json:"role"`
				} `json:"message"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Usage llm.TokenUsage `json:"usage"`
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"base_resp"`
		}

		if err := json.Unmarshal([]byte(data), &raw); err != nil {
			return err
		}
		if raw.Error.Code != 0 && raw.Error.Message != "" {
			return &llm.ProviderError{Code: fmt.Sprintf("%d", raw.Error.Code), Message: raw.Error.Message}
		}
		if len(raw.Choices) == 0 {
			return nil
		}

		content := raw.Choices[0].Delta.Content
		role := raw.Choices[0].Delta.Role
		if content == "" {
			content = raw.Choices[0].Message.Content
		}
		if role == "" {
			role = raw.Choices[0].Message.Role
		}

		cb(llm.ChatResponse{Content: content, Role: llm.Role(role), Usage: raw.Usage})

		if raw.Choices[0].FinishReason != nil && *raw.Choices[0].FinishReason != "" {
			return io.EOF
		}
		return nil
	})
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("minimax stream: %w", err)
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
	var systemPrompt string

	for _, m := range msgs {
		if m.Role == llm.RoleSystem {
			systemPrompt += m.Content + "\n\n"
			continue
		}

		role := "user"
		if m.Role == llm.RoleAI {
			role = "assistant"
		}

		if m.ImageBase64 != "" {
			content := []minimaxContentBlock{
				{Text: m.Content},
				{ImageURL: "data:image/jpeg;base64," + m.ImageBase64},
			}
			result = append(result, map[string]any{
				"role":    role,
				"content": content,
			})
		} else {
			result = append(result, map[string]any{
				"role":    role,
				"content": m.Content,
			})
		}
	}

	if systemPrompt != "" && len(result) > 0 {
		firstMsg := result[0]
		if blocks, ok := firstMsg["content"].([]minimaxContentBlock); ok {
			blocks[0].Text = systemPrompt + blocks[0].Text
			firstMsg["content"] = blocks
		} else if text, ok := firstMsg["content"].(string); ok {
			firstMsg["content"] = systemPrompt + text
		}
	}

	return result
}
