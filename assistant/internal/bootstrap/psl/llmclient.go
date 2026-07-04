package psl

import (
	"context"
	"sync"

	"assistant/pkg/llm"
)

var (
	llmProxySvc   *llm.ProxyService
	llmClient     llm.Client
	onceLLMClient sync.Once
)

func GetProxyService() *llm.ProxyService { return llmProxySvc }

func GetLLMClient() llm.Client { return llmClient }

func InitLLMClient() error {
	onceLLMClient.Do(func() {
		cfg := GetConfig().LLMProxy
		if cfg.VPN == "" {
			cfg.VPN = GetConfig().Settings.VPN
		}
		llmProxySvc = llm.NewProxyService(cfg)
		if llmProxySvc.HasProviders() {
			llmClient = llm.NewProxyClient(llmProxySvc)
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
