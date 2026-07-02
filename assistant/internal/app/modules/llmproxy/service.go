package llmproxy

import (
	"context"
	"net/http"

	"assistant/internal/bootstrap/psl"
	proxysvc "assistant/pkg/llmproxy"
)

type Service struct {
	inner *proxysvc.Service
}

func NewService() *Service {
	cfg := psl.GetConfig().LLMProxy
	return &Service{inner: proxysvc.NewService(cfg)}
}

func (s *Service) Forward(ctx context.Context, reqMap map[string]interface{}, model string) (*http.Response, error) {
	return s.inner.Forward(ctx, reqMap, model)
}

func (s *Service) ActiveProvider() *proxysvc.ProviderConfig {
	return s.inner.ActiveProvider()
}

func (s *Service) ProviderStatuses() []map[string]interface{} {
	return s.inner.ProviderStatuses()
}

func (s *Service) Config() proxysvc.Config {
	return s.inner.Config()
}
