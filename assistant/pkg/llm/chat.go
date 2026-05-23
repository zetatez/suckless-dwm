package llm

import (
	"context"
	"errors"
)

func Chat(ctx context.Context, c Client, req ChatRequest, cb StreamCallback) (*ChatResponse, error) {
	caps := c.Capabilities()

	if cb != nil && caps.Has(CapabilityStream) {
		err := c.StreamChat(ctx, req, cb)
		if err == nil {
			return nil, nil
		}
		if errors.Is(err, ErrNotImplemented) {
		} else {
			return nil, err
		}
	}

	if !caps.Has(CapabilityChat) {
		return nil, ErrNotImplemented
	}

	return c.Chat(ctx, req)
}

func Complete(ctx context.Context, c Client, prompt string, opts ...Option) (*ChatResponse, error) {
	req := ChatRequest{
		Model:       "",
		Messages:    []Message{{Role: RoleUser, Content: prompt}},
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	for _, opt := range opts {
		opt(&req)
	}

	return c.Chat(ctx, req)
}

type Option func(*ChatRequest)

func WithModel(model string) Option {
	return func(r *ChatRequest) { r.Model = model }
}

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
