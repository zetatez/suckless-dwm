package smartapi

import (
	"context"
	"encoding/json"
	"fmt"

	"assistant/pkg/llm"
	"assistant/pkg/utils"
)

type Engine struct {
	client llm.Client
}

func NewEngine(client llm.Client) *Engine {
	return &Engine{client: client}
}

func CompleteJSON[T any](
	ctx context.Context,
	e *Engine,
	prompt string,
	systemPrompt string,
	temperature float32,
	maxTokens int,
) (*T, error) {
	resp, err := llm.Complete(ctx, e.client, prompt, llm.WithSystemPrompt(systemPrompt), llm.WithTemperature(temperature), llm.WithMaxTokens(maxTokens))
	if err != nil {
		return nil, err
	}
	content := utils.CleanJSONResponse(resp.Content)
	var out T
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w\nraw: %s", err, resp.Content)
	}
	return &out, nil
}
