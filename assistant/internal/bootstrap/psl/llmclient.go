package psl

import (
	"context"
	"sync"

	"assistant/pkg/llmproxy"
)

var (
	proxySvc  *llmproxy.Service
	llmClient llmproxy.Client
	onceLLM   sync.Once
)

func GetProxyService() *llmproxy.Service { return proxySvc }

func GetLLMClient() llmproxy.Client { return llmClient }

func InitLLMClient() error {
	onceLLM.Do(func() {
		cfg := GetConfig().LLMProxy
		proxySvc = llmproxy.NewService(cfg)
		if proxySvc.HasProviders() {
			llmClient = llmproxy.NewProxyClient(proxySvc)
		}
	})
	return nil
}

func RegisterCleanupLLM() {
	RegisterCleanup(func(ctx context.Context) {
		llmClient = nil
		proxySvc = nil
	})
}
