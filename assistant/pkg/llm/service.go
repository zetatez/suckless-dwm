package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
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
	NeedVPN  bool     `mapstructure:"need_vpn"`
}

type Config struct {
	ProxiedModel  string           `mapstructure:"proxied_model"`
	ProbeInterval int              `mapstructure:"probe_interval"`
	ProxiedAPIKey string           `mapstructure:"proxied_api_key"`
	Timeout       int              `mapstructure:"timeout"`
	Temperature   float32          `mapstructure:"temperature"`
	VPN           string           `mapstructure:"vpn"`
	Providers     []ProviderConfig `mapstructure:"providers"`
}

type providerState struct {
	ProviderConfig
	status           ProviderStatus
	rateLimitedUntil time.Time
}

type ProxyService struct {
	config       Config
	mu           sync.RWMutex
	providers    []*providerState
	active       *providerState
	lastModel    string
	httpClient   *http.Client
	vpnClient    *http.Client
	lastActivity time.Time
	probeMu      sync.Mutex
	probeRunning bool
}

func NewProxyService(cfg Config) *ProxyService {
	timeout := 5 * time.Minute
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	s := &ProxyService{
		config: cfg,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: sharedTransport,
		},
	}

	if cfg.VPN != "" {
		if dialer, err := newSocksDialer(cfg.VPN); err == nil {
			tr := &http.Transport{
				MaxIdleConns:        8,
				MaxIdleConnsPerHost: 2,
				IdleConnTimeout:     120 * time.Second,
				DisableCompression:  false,
			}
			tr.DialContext = dialer.DialContext
			s.vpnClient = &http.Client{Timeout: timeout, Transport: tr}
		}
	}

	for _, pc := range cfg.Providers {
		s.providers = append(s.providers, &providerState{
			ProviderConfig: pc,
			status:         StatusAvailable,
		})
	}
	return s
}

func newSocksDialer(proxyURL string) (interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	d, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &socksDialer{d: d}, nil
}

type socksDialer struct {
	d proxy.Dialer
}

func (s *socksDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return s.d.Dial(network, addr)
}

func (s *ProxyService) Config() Config { return s.config }

func (s *ProxyService) HasProviders() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.providers) > 0
}

func (s *ProxyService) ActiveProvider() *ProviderConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.active == nil {
		return nil
	}
	return &s.active.ProviderConfig
}

func (s *ProxyService) ProviderStatuses() []map[string]interface{} {
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

func (s *ProxyService) pickBestLocked(model string) *providerState {
	now := time.Now()
	for _, p := range s.providers {
		if p.status != StatusAvailable && p.PlanType != "payg" {
			continue
		}
		if p.status == StatusExhausted {
			continue
		}
		if p.status == StatusAvailable && now.Before(p.rateLimitedUntil) {
			continue
		}
		if model != "" && !slices.Contains(p.Models, model) {
			continue
		}
		return p
	}
	return nil
}

func (s *ProxyService) markProvider(p *providerState, status ProviderStatus, cooldown time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if status != "" {
		p.status = status
	}
	if cooldown > 0 {
		p.rateLimitedUntil = time.Now().Add(cooldown)
	}
	if s.active == p {
		s.active = nil
	}
}

func (s *ProxyService) ProbeHigherPriority() {
	s.mu.RLock()
	active := s.active
	var candidates []*providerState
	if active == nil {
		for _, p := range s.providers {
			if p.status != StatusAvailable {
				candidates = append(candidates, p)
			}
		}
	} else if len(s.providers) > 0 && active != s.providers[0] {
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
		if p.PlanType == "payg" {
			continue
		}
		client := s.httpClient
		if p.NeedVPN && s.vpnClient != nil {
			client = s.vpnClient
		}
		if probeProvider(p, client) {
			s.mu.Lock()
			p.status = StatusAvailable
			p.rateLimitedUntil = time.Time{}
			s.mu.Unlock()
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	model := s.lastModel
	if s.active == nil {
		if best := s.pickBestLocked(model); best != nil {
			s.active = best
			log.Printf("[llmproxy] probe recovered, switched to provider: %s", best.Name)
		}
		return
	}
	if best := s.pickBestLocked(model); best != nil && best != s.active {
		s.active = best
		log.Printf("[llmproxy] probe recovered, switched to provider: %s", best.Name)
	}
}

func probeProvider(p *providerState, client *http.Client) bool {
	if len(p.Models) == 0 {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	probeModel := p.Models[0]
	payload, _ := json.Marshal(map[string]interface{}{
		"model":      probeModel,
		"messages":   []map[string]string{{"role": "user", "content": "hi"}},
		"max_tokens": 1,
		"stream":     false,
	})
	req, err := http.NewRequestWithContext(ctx, "POST",
		strings.TrimRight(p.BaseURL, "/")+"/chat/completions",
		bytes.NewReader(payload))
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "assistant/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode == http.StatusOK
}

func (s *ProxyService) NotifyActivity() {
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

func (s *ProxyService) probeLoop() {
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

func (s *ProxyService) ForwardChat(ctx context.Context, cfg ProviderConfig, bodyReader io.ReadCloser) (*http.Response, error) {
	targetURL := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "assistant/1.0")

	client := s.httpClient
	if cfg.NeedVPN && s.vpnClient != nil {
		client = s.vpnClient
	}
	return client.Do(req)
}

func (s *ProxyService) findNext(current *providerState, modelSpecific bool, originalModel string) *providerState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	for _, p := range s.providers {
		if p == current {
			continue
		}
		if p.status != StatusAvailable && p.PlanType != "payg" {
			continue
		}
		if p.status == StatusExhausted {
			continue
		}
		if p.status == StatusAvailable && now.Before(p.rateLimitedUntil) {
			continue
		}
		if modelSpecific && !slices.Contains(p.Models, originalModel) {
			continue
		}
		return p
	}
	return nil
}

func (s *ProxyService) Forward(ctx context.Context, reqMap map[string]interface{}, requestedModel string) (*http.Response, error) {
	var p *providerState
	var modelSpecific bool

	s.mu.Lock()
	if requestedModel == s.config.ProxiedModel || requestedModel == "" {
		s.lastModel = ""
		p = s.pickBestLocked("")
	} else {
		s.lastModel = requestedModel
		p = s.pickBestLocked(requestedModel)
		if p == nil {
			p = s.pickBestLocked("")
		} else {
			modelSpecific = true
		}
	}
	if p != nil {
		s.active = p
		log.Printf("[llmproxy] selected provider: %s (model=%s)", p.Name, resolveModel(p, requestedModel))
	}
	s.mu.Unlock()

	if p == nil {
		s.NotifyActivity()
		return nil, s.noProviderError()
	}

	for {
		reqMap["model"] = resolveModel(p, requestedModel)
		stripKnownBad(reqMap)
		body, _ := json.Marshal(reqMap)

		resp, err := s.ForwardChat(ctx, p.ProviderConfig, io.NopCloser(bytes.NewReader(body)))
		if err != nil {
			logError("forward to %s failed: %v", p.Name, err)
			s.markProvider(p, StatusOffline, 0)
			next := s.findNext(p, modelSpecific, requestedModel)
			if next == nil {
				s.NotifyActivity()
				return nil, fmt.Errorf("all providers failed: %s: %v", p.Name, err)
			}
			p = next
			log.Printf("[llmproxy] failover to provider: %s", p.Name)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			s.NotifyActivity()
			return resp, nil
		}

		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		resp.Body.Close()

		if resp.StatusCode == 429 {
			s.markProvider(p, "", 60*time.Second)
			log.Printf("[llmproxy] provider %s rate limited (429), cool down 60s", p.Name)
		} else if p.PlanType == "fixed" {
			s.markProvider(p, StatusExhausted, 0)
			log.Printf("[llmproxy] provider %s error: %s (code=%d)", p.Name, providerErrorMessage(bodyBytes), resp.StatusCode)
		} else {
			s.markProvider(p, StatusOffline, 0)
			log.Printf("[llmproxy] provider %s error: %s (code=%d)", p.Name, providerErrorMessage(bodyBytes), resp.StatusCode)
		}

		next := s.findNext(p, modelSpecific, requestedModel)
		if next == nil {
			s.NotifyActivity()
			errMsg := providerErrorMessage(bodyBytes)
			return nil, &HTTPError{Code: resp.StatusCode, Message: errMsg}
		}
		p = next
		log.Printf("[llmproxy] failover to provider: %s", p.Name)
	}
}

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string { return e.Message }

var knownBadFields = []string{"promptCacheKey"}

func stripKnownBad(req map[string]interface{}) {
	for _, k := range knownBadFields {
		delete(req, k)
	}
}

func resolveModel(p *providerState, requested string) string {
	if len(p.Models) == 0 {
		return requested
	}
	if requested != "" && slices.Contains(p.Models, requested) {
		return requested
	}
	return p.Models[0]
}

func logError(format string, args ...interface{}) {
	log.Printf("[llmproxy] "+format, args...)
}

func (s *ProxyService) noProviderError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exhausted, offline := 0, 0
	for _, p := range s.providers {
		if p.status == StatusExhausted {
			exhausted++
		} else if p.status == StatusOffline && p.PlanType != "payg" {
			offline++
		}
	}
	if exhausted > 0 {
		return fmt.Errorf("no available provider (%d exhausted, %d offline)", exhausted, offline)
	}
	return fmt.Errorf("no available provider (%d offline)", offline)
}

func providerErrorMessage(body []byte) string {
	if len(body) == 0 {
		return "provider failed"
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return string(body)
	}
	if msg, ok := parsed["error"].(string); ok {
		return msg
	}
	if errObj, ok := parsed["error"].(map[string]interface{}); ok {
		if msg, ok := errObj["message"].(string); ok {
			return msg
		}
	}
	return string(body)
}
