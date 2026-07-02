package llmproxy

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"assistant/internal/bootstrap/psl"
)

var sharedTransport = &http.Transport{
	MaxIdleConns:        8,
	MaxIdleConnsPerHost: 2,
	IdleConnTimeout:     120 * time.Second,
	DisableCompression:  false,
}

type ProviderStatus string

const (
	StatusAvailable ProviderStatus = "available"
	StatusExhausted ProviderStatus = "exhausted"
	StatusOffline   ProviderStatus = "offline"
)

type providerState struct {
	psl.ProviderConfig
	status ProviderStatus
}

type Service struct {
	mu           sync.RWMutex
	providers    []*providerState
	active       *providerState
	httpClient   *http.Client
	lastActivity time.Time
	probeMu      sync.Mutex
	probeRunning bool
}

func NewService() *Service {
	s := &Service{
		httpClient: &http.Client{
			Timeout:   5 * time.Minute,
			Transport: sharedTransport,
		},
	}
	for _, pc := range psl.GetConfig().LLMProxy.Providers {
		s.providers = append(s.providers, &providerState{
			ProviderConfig: pc,
			status:         StatusAvailable,
		})
	}
	return s
}

func (s *Service) ActiveProvider() *psl.ProviderConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.active == nil {
		return nil
	}
	return &s.active.ProviderConfig
}

func (s *Service) ProviderStatuses() []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []map[string]interface{}
	for _, p := range s.providers {
		m := map[string]interface{}{
			"name":      p.Name,
			"status":    p.status,
			"plan_type": p.PlanType,
		}
		res = append(res, m)
	}
	return res
}

func (s *Service) selectProvider() *providerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.pickBestLocked()
	if p != nil {
		s.active = p
	}
	return p
}

func (s *Service) selectProviderByModel(model string) *providerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	var fixed, payg []*providerState
	for _, p := range s.providers {
		if p.status != StatusAvailable {
			continue
		}
		if !hasModel(p.Models, model) {
			continue
		}
		if p.PlanType == "fixed" {
			fixed = append(fixed, p)
		} else {
			payg = append(payg, p)
		}
	}
	if len(fixed) > 0 {
		s.active = fixed[0]
		return fixed[0]
	}
	if len(payg) > 0 {
		s.active = payg[0]
		return payg[0]
	}
	return nil
}

func (s *Service) pickBestLocked() *providerState {
	var fixed, payg []*providerState
	for _, p := range s.providers {
		if p.status != StatusAvailable {
			continue
		}
		if p.PlanType == "fixed" {
			fixed = append(fixed, p)
		} else {
			payg = append(payg, p)
		}
	}
	if len(fixed) > 0 {
		return fixed[0]
	}
	if len(payg) > 0 {
		return payg[0]
	}
	return nil
}

func (s *Service) MarkExhausted(p *providerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.status = StatusExhausted
	if s.active == p {
		s.active = nil
	}
}

func (s *Service) MarkOffline(p *providerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.status = StatusOffline
	if s.active == p {
		s.active = nil
	}
}

// ProbeHigherPriority 探测当前 active 前面的高优先级供应商。
// 如果发现可用的高优先级供应商，切换到它。由 handler 在每次请求完成后调用。
func (s *Service) ProbeHigherPriority() {
	s.mu.Lock()
	defer s.mu.Unlock()

	active := s.active
	if active == nil {
		if best := s.pickBestLocked(); best != nil {
			s.active = best
		}
		return
	}
	if active == s.providers[0] {
		return
	}

	for _, p := range s.providers {
		if p == active {
			break
		}
		if p.status == StatusAvailable {
			continue
		}
		if probeProvider(p, s.httpClient) {
			p.status = StatusAvailable
		}
	}

	if best := s.pickBestLocked(); best != nil && best != active {
		s.active = best
	}
}

func probeProvider(p *providerState, client *http.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET",
		strings.TrimRight(p.BaseURL, "/")+"/models", nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("User-Agent", "assistant/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode == http.StatusOK
}

// NotifyActivity 通知路由器有请求完成，触发探测循环（如有需要）。
func (s *Service) NotifyActivity() {
	s.probeMu.Lock()
	s.lastActivity = time.Now()
	if !s.probeRunning {
		s.probeRunning = true
		s.probeMu.Unlock()
		go s.probeLoop()
	} else {
		s.probeMu.Unlock()
	}
}

func (s *Service) probeLoop() {
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.ProbeHigherPriority()

		s.probeMu.Lock()
		if time.Since(s.lastActivity) > interval*2 {
			s.probeRunning = false
			s.probeMu.Unlock()
			return
		}
		s.probeMu.Unlock()
	}
}

func (s *Service) ForwardChat(ctx context.Context, cfg psl.ProviderConfig, bodyReader io.ReadCloser) (*http.Response, error) {
	targetURL := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "assistant/1.0")

	return s.httpClient.Do(req)
}

func logError(format string, args ...interface{}) {
	log.Printf("[llmproxy] "+format, args...)
}

func hasModel(models []string, target string) bool {
	for _, m := range models {
		if m == target {
			return true
		}
	}
	return false
}
