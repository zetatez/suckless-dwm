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
	MiddleModel   string           `mapstructure:"middle_model"`
	ProbeInterval int              `mapstructure:"probe_interval"`
	AuthToken     string           `mapstructure:"auth_token"`
	Timeout       int              `mapstructure:"timeout"`
	Temperature   float32          `mapstructure:"temperature"`
	VPNProxy      string           `mapstructure:"vpn"`
	Providers     []ProviderConfig `mapstructure:"providers"`
}

type providerState struct {
	ProviderConfig
	status ProviderStatus
}

type ProxyService struct {
	config       Config
	mu           sync.RWMutex
	providers    []*providerState
	active       *providerState
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

	if cfg.VPNProxy != "" {
		if dialer, err := newSocksDialer(cfg.VPNProxy); err == nil {
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

func (s *ProxyService) selectProvider(model string) *providerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.pickBestLocked(model)
	if p != nil {
		s.active = p
	}
	return p
}

func (s *ProxyService) pickBestLocked(model string) *providerState {
	var fixed, payg []*providerState
	for _, p := range s.providers {
		if p.status != StatusAvailable && p.PlanType != "payg" {
			continue
		}
		if p.status == StatusExhausted {
			continue
		}
		if model != "" && !hasModel(p.Models, model) {
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

func (s *ProxyService) markStatus(p *providerState, status ProviderStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.status = status
	if s.active == p {
		s.active = nil
	}
}

func (s *ProxyService) ProbeHigherPriority() {
	s.mu.RLock()
	active := s.active
	var candidates []*providerState
	if active == nil {
		// 全部不可用，探测所有非 available 供应商
		for _, p := range s.providers {
			if p.status != StatusAvailable {
				candidates = append(candidates, p)
			}
		}
	} else if active != s.providers[0] {
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
		client := s.httpClient
		if p.NeedVPN && s.vpnClient != nil {
			client = s.vpnClient
		}
		if probeProvider(p, client) {
			s.mu.Lock()
			p.status = StatusAvailable
			s.mu.Unlock()
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active == nil {
		if best := s.pickBestLocked(""); best != nil {
			s.active = best
		}
		return
	}
	if best := s.pickBestLocked(""); best != nil && best != s.active {
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
		if modelSpecific && !hasModel(p.Models, originalModel) {
			continue
		}
		return p
	}
	return nil
}

func (s *ProxyService) Forward(ctx context.Context, reqMap map[string]interface{}, requestedModel string) (*http.Response, error) {
	var p *providerState
	var modelSpecific bool

	if requestedModel == s.config.MiddleModel || requestedModel == "" {
		p = s.selectProvider("")
	} else {
		p = s.selectProvider(requestedModel)
		if p == nil {
			p = s.selectProvider("")
		} else {
			modelSpecific = true
		}
	}

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
			s.markStatus(p, StatusOffline)
			next := s.findNext(p, modelSpecific, requestedModel)
			if next == nil {
				s.NotifyActivity()
				return nil, fmt.Errorf("all providers failed: %s: %v", p.Name, err)
			}
			p = next
			continue
		}

		if resp.StatusCode == http.StatusOK {
			s.NotifyActivity()
			return resp, nil
		}

		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		resp.Body.Close()

		if p.PlanType == "fixed" {
			s.markStatus(p, StatusExhausted)
		} else {
			s.markStatus(p, StatusOffline)
		}

		next := s.findNext(p, modelSpecific, requestedModel)
		if next == nil {
			s.NotifyActivity()
			errMsg := providerErrorMessage(bodyBytes)
			return nil, &HTTPError{Code: resp.StatusCode, Message: errMsg}
		}
		p = next
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
