package psl

import (
	"context"
	"sync"

	"assistant/pkg/llmproxy"
)

var (
	llmProxySvc   *llmproxy.Service
	llmClient     llmproxy.Client
	onceLLMClient sync.Once
)

func GetProxyService() *llmproxy.Service { return llmProxySvc }

func GetLLMClient() llmproxy.Client { return llmClient }

func InitLLMClient() error {
	onceLLMClient.Do(func() {
		cfg := GetConfig().LLMProxy
		llmProxySvc = llmproxy.NewService(cfg)
		if llmProxySvc.HasProviders() {
			llmClient = llmproxy.NewProxyClient(llmProxySvc)
		}
	})
	return nil
}

func RegisterCleanupLLM() {
	RegisterCleanup(func(ctx context.Context) {
		llmClient = nil
		llmProxySvc = nil
	})
}
