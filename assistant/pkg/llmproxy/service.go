package llmproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ErrNoProvider         = errors.New("no available provider")
	ErrAllProvidersFailed = errors.New("all providers failed")
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

type ProviderConfig struct {
	Name     string   `mapstructure:"name"`
	BaseURL  string   `mapstructure:"base_url"`
	APIKey   string   `mapstructure:"api_key"`
	Models   []string `mapstructure:"models"`
	PlanType string   `mapstructure:"plan_type"`
}

type Config struct {
	MiddleModel   string           `mapstructure:"middle_model"`
	ProbeInterval int              `mapstructure:"probe_interval"`
	AuthToken     string           `mapstructure:"auth_token"`
	Timeout       int              `mapstructure:"timeout"`
	Temperature   float32          `mapstructure:"temperature"`
	Providers     []ProviderConfig `mapstructure:"providers"`
}

type providerState struct {
	ProviderConfig
	status ProviderStatus
}

type Service struct {
	config       Config
	mu           sync.RWMutex
	providers    []*providerState
	active       *providerState
	httpClient   *http.Client
	lastActivity time.Time
	probeMu      sync.Mutex
	probeRunning bool
}

func NewService(cfg Config) *Service {
	timeout := 5 * time.Minute
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}
	s := &Service{
		config: cfg,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: sharedTransport,
		},
	}
	for _, pc := range cfg.Providers {
		s.providers = append(s.providers, &providerState{
			ProviderConfig: pc,
			status:         StatusAvailable,
		})
	}
	return s
}

func (s *Service) Config() Config { return s.config }

func (s *Service) HasProviders() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.providers) > 0
}

func (s *Service) ActiveProvider() *ProviderConfig {
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

func (s *Service) markExhausted(p *providerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.status = StatusExhausted
	if s.active == p {
		s.active = nil
	}
}

func (s *Service) markOffline(p *providerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.status = StatusOffline
	if s.active == p {
		s.active = nil
	}
}

func (s *Service) ProbeHigherPriority() {
	s.mu.RLock()
	active := s.active
	var candidates []*providerState
	if active != nil && active != s.providers[0] {
		for _, p := range s.providers {
			if p == active {
				break
			}
			if p.status != StatusAvailable {
				candidates = append(candidates, p)
			}
		}
	}
	s.mu.RUnlock()

	for _, p := range candidates {
		if probeProvider(p, s.httpClient) {
			s.mu.Lock()
			p.status = StatusAvailable
			s.mu.Unlock()
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active == nil {
		if best := s.pickBestLocked(); best != nil {
			s.active = best
		}
		return
	}
	if best := s.pickBestLocked(); best != nil && best != s.active {
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
	interval := 30 * time.Second
	if s.config.ProbeInterval > 0 {
		interval = time.Duration(s.config.ProbeInterval) * time.Second
	}
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

func (s *Service) ForwardChat(ctx context.Context, cfg ProviderConfig, bodyReader io.ReadCloser) (*http.Response, error) {
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

func (s *Service) findNext(current *providerState, modelSpecific bool, originalModel string) *providerState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.providers {
		if p == current {
			continue
		}
		if p.status != StatusAvailable {
			continue
		}
		if modelSpecific && !hasModel(p.Models, originalModel) {
			continue
		}
		return p
	}
	return nil
}

func (s *Service) Forward(ctx context.Context, reqMap map[string]interface{}, requestedModel string) (*http.Response, error) {
	var p *providerState
	var modelSpecific bool

	if requestedModel == s.config.MiddleModel || requestedModel == "" {
		p = s.selectProvider()
	} else {
		p = s.selectProviderByModel(requestedModel)
		if p == nil {
			p = s.selectProvider()
		} else {
			modelSpecific = true
		}
	}

	if p == nil {
		return nil, ErrNoProvider
	}

	for {
		reqMap["model"] = resolveModel(p, requestedModel)
		body, _ := json.Marshal(reqMap)

		resp, err := s.ForwardChat(ctx, p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			logError("forward to %s failed: %v", p.Name, err)
			s.markOffline(p)
			next := s.findNext(p, modelSpecific, requestedModel)
			if next == nil {
				return nil, ErrAllProvidersFailed
			}
			p = next
			continue
		}

		if resp.StatusCode == http.StatusOK {
			s.NotifyActivity()
			return resp, nil
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if p.PlanType == "fixed" {
			s.markExhausted(p)
		} else {
			s.markOffline(p)
		}

		next := s.findNext(p, modelSpecific, requestedModel)
		if next == nil {
			return nil, &HTTPError{Code: resp.StatusCode, Message: "provider failed"}
		}
		p = next
	}
}

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string { return e.Message }

func resolveModel(p *providerState, requested string) string {
	if requested == "" {
		return p.Models[0]
	}
	if hasModel(p.Models, requested) {
		return requested
	}
	return p.Models[0]
}

func hasModel(models []string, target string) bool {
	for _, m := range models {
		if m == target {
			return true
		}
	}
	return false
}

func logError(format string, args ...interface{}) {
	log.Printf("[llmproxy] "+format, args...)
}
