package llm

import "context"

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleAI     Role = "assistant"
)

type Message struct {
	Role        Role       `json:"role"`
	Content     string     `json:"content"`
	ToolCalls   []ToolCall `json:"tool_calls,omitempty"`
	ImageBase64 string     `json:"-"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float32
	MaxTokens   int
}

type ChatResponse struct {
	Content   string
	Role      Role
	ToolCalls []ToolCall
	Usage     TokenUsage
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Client interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

type Option func(*ChatRequest)

func WithTemperature(temp float32) Option {
	return func(r *ChatRequest) { r.Temperature = temp }
}

func WithMaxTokens(max int) Option {
	return func(r *ChatRequest) { r.MaxTokens = max }
}

func WithSystemPrompt(system string) Option {
	return func(r *ChatRequest) {
		r.Messages = append([]Message{{Role: RoleSystem, Content: system}}, r.Messages...)
	}
}

func WithImageBase64(img string) Option {
	return func(r *ChatRequest) {
		if len(r.Messages) > 0 {
			r.Messages[len(r.Messages)-1].ImageBase64 = img
		}
	}
}

func Complete(ctx context.Context, c Client, prompt string, opts ...Option) (*ChatResponse, error) {
	req := ChatRequest{
		Messages:    []Message{{Role: RoleUser, Content: prompt}},
		Temperature: 0.7,
		MaxTokens:   4096,
	}
	for _, opt := range opts {
		opt(&req)
	}
	return c.Chat(ctx, req)
}
