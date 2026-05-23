package gemini

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"assistant/pkg/llm"
)

// Gemini provider implementation using the Google Generative Language API.
// Docs: https://ai.google.dev/api/rest

type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *llm.BaseClient
	http    *http.Client
}

func init() {
	llm.Register("gemini", New)
}

func New(cfg llm.Config) (llm.Client, error) {
	if cfg.APIKey == "" {
		return nil, llm.ErrInvalidConfig
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		// Keep version in baseURL so paths stay small.
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	model := cfg.Model
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		client:  llm.NewBaseClient(baseURL, cfg),
		http:    cfg.GetHTTPClient(),
	}, nil
}

func (c *Client) Provider() string { return "gemini" }

func (c *Client) Model() string { return c.model }

func (c *Client) Capabilities() llm.Capabilities {
	return llm.Capabilities{Supported: llm.CapabilityChat | llm.CapabilityStream}
}

func (c *Client) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	payload := c.buildPayload(req)

	headers := map[string]string{
		"x-goog-api-key": c.apiKey,
		"Content-Type":   "application/json",
	}

	path := fmt.Sprintf("/models/%s:generateContent", c.getModel(req.Model))
	resp, err := c.client.Do(ctx, "POST", path, payload, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, parseGeminiError(resp.Body, resp.StatusCode)
	}

	var raw geminiGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Candidates) == 0 {
		return nil, llm.ErrMaxRetries
	}

	content, toolCalls := extractTextAndToolCalls(raw.Candidates[0].Content)
	role := llm.RoleAI
	if raw.Candidates[0].Content.Role == "user" {
		role = llm.RoleUser
	}

	return &llm.ChatResponse{
		Content:   content,
		Role:      role,
		ToolCalls: toolCalls,
		Usage: llm.TokenUsage{
			PromptTokens:     raw.UsageMetadata.PromptTokenCount,
			CompletionTokens: raw.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      raw.UsageMetadata.TotalTokenCount,
		},
		Raw: raw,
	}, nil
}

func (c *Client) StreamChat(ctx context.Context, req llm.ChatRequest, cb llm.StreamCallback) error {
	payload := c.buildPayload(req)

	// SSE endpoint.
	path := fmt.Sprintf("/models/%s:streamGenerateContent?alt=sse", c.getModel(req.Model))

	b, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(b))
	httpReq.Header.Set("x-goog-api-key", c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return parseGeminiError(resp.Body, resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}

		var raw geminiGenerateResponse
		if err := json.Unmarshal([]byte(data), &raw); err != nil {
			// Ignore non-JSON events.
			continue
		}
		if len(raw.Candidates) == 0 {
			continue
		}

		text, toolCalls := extractTextAndToolCalls(raw.Candidates[0].Content)
		if text == "" && len(toolCalls) == 0 {
			continue
		}

		cb(llm.ChatResponse{
			Content:   text,
			Role:      llm.RoleAI,
			ToolCalls: toolCalls,
			Usage: llm.TokenUsage{
				PromptTokens:     raw.UsageMetadata.PromptTokenCount,
				CompletionTokens: raw.UsageMetadata.CandidatesTokenCount,
				TotalTokens:      raw.UsageMetadata.TotalTokenCount,
			},
			Raw: raw,
		})
	}

	return scanner.Err()
}

func (c *Client) getModel(model string) string {
	if model != "" {
		return model
	}
	return c.model
}

func (c *Client) buildPayload(req llm.ChatRequest) map[string]any {
	system, contents := toGeminiContents(req.Messages)
	if len(contents) == 0 {
		contents = []map[string]any{{
			"role":  "user",
			"parts": []map[string]any{{"text": ""}},
		}}
	}

	payload := map[string]any{
		"contents": contents,
	}

	if system != "" {
		payload["systemInstruction"] = map[string]any{
			"parts": []map[string]any{{"text": system}},
		}
	}

	gen := map[string]any{}
	if req.Temperature > 0 {
		gen["temperature"] = req.Temperature
	}
	if req.TopP > 0 {
		gen["topP"] = req.TopP
	}
	if req.MaxTokens > 0 {
		gen["maxOutputTokens"] = req.MaxTokens
	}
	if len(gen) > 0 {
		payload["generationConfig"] = gen
	}

	return payload
}

func toGeminiContents(msgs []llm.Message) (system string, contents []map[string]any) {
	var sysParts []string
	for _, m := range msgs {
		switch m.Role {
		case llm.RoleSystem:
			if strings.TrimSpace(m.Content) != "" {
				sysParts = append(sysParts, m.Content)
			}
		case llm.RoleUser, llm.RoleAI:
			if strings.TrimSpace(m.Content) == "" {
				continue
			}
			role := "user"
			if m.Role == llm.RoleAI {
				role = "model"
			}
			contents = append(contents, map[string]any{
				"role":  role,
				"parts": []map[string]any{{"text": m.Content}},
			})
		default:
			// Ignore unknown roles.
		}
	}
	return strings.Join(sysParts, "\n"), contents
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Role  string       `json:"role"`
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

type geminiPart struct {
	Text         string `json:"text,omitempty"`
	FunctionCall *struct {
		Name string         `json:"name"`
		Args map[string]any `json:"args"`
	} `json:"functionCall,omitempty"`
}

func extractTextAndToolCalls(content struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}) (string, []llm.ToolCall) {
	var sb strings.Builder
	var toolCalls []llm.ToolCall
	for i, p := range content.Parts {
		if p.Text != "" {
			sb.WriteString(p.Text)
		}
		if p.FunctionCall != nil {
			argsBytes, _ := json.Marshal(p.FunctionCall.Args)
			toolCalls = append(toolCalls, llm.ToolCall{
				ID:   fmt.Sprintf("call_%d", i),
				Type: "function",
				Function: llm.Function{
					Name:      p.FunctionCall.Name,
					Arguments: string(argsBytes),
				},
			})
		}
	}
	return sb.String(), toolCalls
}

func parseGeminiError(r io.Reader, statusCode int) error {
	var raw struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}
	b, _ := io.ReadAll(r)
	if err := json.Unmarshal(b, &raw); err != nil {
		return &llm.ProviderError{Code: fmt.Sprintf("http_%d", statusCode), Message: string(bytes.TrimSpace(b))}
	}
	code := raw.Error.Status
	if code == "" {
		code = fmt.Sprintf("http_%d", statusCode)
	}
	msg := raw.Error.Message
	if msg == "" {
		msg = fmt.Sprintf("gemini error (http %d)", statusCode)
	}
	return &llm.ProviderError{Code: code, Message: msg, Raw: json.RawMessage(b)}
}
