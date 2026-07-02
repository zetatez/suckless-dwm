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

type Router struct {
	mu           sync.RWMutex
	providers    []*providerState
	active       *providerState
	httpClient   *http.Client
	lastActivity time.Time
	probeMu      sync.Mutex
	probeRunning bool
}

var sharedTransport = &http.Transport{
	MaxIdleConns:        8,
	MaxIdleConnsPerHost: 2,
	IdleConnTimeout:     120 * time.Second,
	DisableCompression:  false,
}

func NewRouter() *Router {
	r := &Router{
		httpClient: &http.Client{
			Timeout:   5 * time.Minute,
			Transport: sharedTransport,
		},
	}
	for _, pc := range psl.GetConfig().LLMProxy.Providers {
		r.providers = append(r.providers, &providerState{
			ProviderConfig: pc,
			status:         StatusAvailable,
		})
	}
	return r
}

func (r *Router) ActiveProvider() *psl.ProviderConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.active == nil {
		return nil
	}
	return &r.active.ProviderConfig
}

func (r *Router) ProviderStatuses() []map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []map[string]interface{}
	for _, p := range r.providers {
		m := map[string]interface{}{
			"name":      p.Name,
			"status":    p.status,
			"plan_type": p.PlanType,
		}
		res = append(res, m)
	}
	return res
}

func (r *Router) selectProvider() *providerState {
	r.mu.Lock()
	defer r.mu.Unlock()
	p := r.pickBestLocked()
	if p != nil {
		r.active = p
	}
	return p
}

func (r *Router) selectProviderByModel(model string) *providerState {
	r.mu.Lock()
	defer r.mu.Unlock()
	var fixed, payg []*providerState
	for _, p := range r.providers {
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
		r.active = fixed[0]
		return fixed[0]
	}
	if len(payg) > 0 {
		r.active = payg[0]
		return payg[0]
	}
	return nil
}

func (r *Router) pickBestLocked() *providerState {
	var fixed, payg []*providerState
	for _, p := range r.providers {
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

func (r *Router) MarkExhausted(p *providerState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.status = StatusExhausted
	if r.active == p {
		r.active = nil
	}
}

func (r *Router) MarkOffline(p *providerState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.status = StatusOffline
	if r.active == p {
		r.active = nil
	}
}

// ProbeHigherPriority 探测当前 active 前面的高优先级供应商。
// 如果发现可用的高优先级供应商，切换到它。由 handler 在每次请求完成后调用。
func (r *Router) ProbeHigherPriority() {
	r.mu.Lock()
	defer r.mu.Unlock()

	active := r.active
	if active == nil {
		if best := r.pickBestLocked(); best != nil {
			r.active = best
		}
		return
	}
	if active == r.providers[0] {
		return
	}

	for _, p := range r.providers {
		if p == active {
			break
		}
		if p.status == StatusAvailable {
			continue
		}
		if probeProvider(p, r.httpClient) {
			p.status = StatusAvailable
		}
	}

	if best := r.pickBestLocked(); best != nil && best != active {
		r.active = best
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
func (r *Router) NotifyActivity() {
	r.probeMu.Lock()
	r.lastActivity = time.Now()
	if !r.probeRunning {
		r.probeRunning = true
		r.probeMu.Unlock()
		go r.probeLoop()
	} else {
		r.probeMu.Unlock()
	}
}

func (r *Router) probeLoop() {
	const interval = 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		r.ProbeHigherPriority()

		r.probeMu.Lock()
		if time.Since(r.lastActivity) > interval*2 {
			r.probeRunning = false
			r.probeMu.Unlock()
			return
		}
		r.probeMu.Unlock()
	}
}

func (r *Router) ForwardChat(ctx context.Context, cfg psl.ProviderConfig, bodyReader io.ReadCloser) (*http.Response, error) {
	targetURL := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "assistant/1.0")

	return r.httpClient.Do(req)
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
