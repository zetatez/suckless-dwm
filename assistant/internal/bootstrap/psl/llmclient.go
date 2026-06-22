package psl

import (
	"context"
	"fmt"
	"sync"

	"assistant/pkg/llm"
)

var (
	llmClient llm.Client
	onceLLM   sync.Once
	llmErr    error
)

func GetLLMClient() llm.Client { return llmClient }

func InitLLMClient() error {
	cfg := GetConfig().LLM
	onceLLM.Do(func() {
		llmClient, llmErr = llm.NewClient(cfg.Provider, llm.Config{
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			Extra:       cfg.Extra,
			Timeout:     cfg.Timeout,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
		})
		if llmErr != nil {
			llmErr = fmt.Errorf("create LLM client: %w", llmErr)
		}
	})
	return llmErr
}

func RegisterCleanupLLM() {
	RegisterCleanup(func(ctx context.Context) {
		if llmClient != nil {
			llmClient = nil
		}
	})
}
