package llm

import (
	"context"
)

type Client interface {
	Provider() string
	Model() string
	Capabilities() Capabilities
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	StreamChat(ctx context.Context, req ChatRequest, cb StreamCallback) error
}

type Embedder interface {
	Client
	Embed(ctx context.Context, texts []string) ([]Embedding, error)
}

type Embedding struct {
	Vector []float32
	Model  string
}
