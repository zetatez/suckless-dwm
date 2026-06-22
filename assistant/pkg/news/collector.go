package news

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"assistant/pkg/dwmblocknotify"
)

type Item struct {
	Title  string
	Source string
	Link   string
}

type Provider interface {
	Name() string
	Fetch(ctx context.Context, limit int) ([]Item, error)
}

type Collector struct {
	providers map[string]Provider
	client    *http.Client
}

func New() *Collector {
	client := &http.Client{Timeout: 8 * time.Second}
	c := &Collector{
		providers: make(map[string]Provider),
		client:    client,
	}
	c.Register(&rssProvider{name: "top-news", source: "新闻", url: topNewsURL, client: client, cacheTTL: 3 * time.Minute})
	c.Register(&rssProvider{name: "finance-news", source: "财经", url: financeNewsURL, client: client, cacheTTL: 3 * time.Minute})
	return c
}

func (c *Collector) Register(p Provider) {
	c.providers[p.Name()] = p
}

func (c *Collector) Providers() []string {
	providers := make([]string, 0, len(c.providers))
	for name := range c.providers {
		providers = append(providers, name)
	}
	sort.Strings(providers)
	return providers
}

func (c *Collector) Fetch(ctx context.Context, provider string, limit int) ([]Item, error) {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = "top-news"
	}
	p, ok := c.providers[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
	if limit <= 0 {
		limit = 5
	}
	items, err := p.Fetch(ctx, limit)
	if err != nil {
		return nil, err
	}
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (c *Collector) Notify(ctx context.Context, provider string, limit int, ttl, interval time.Duration) error {
	items, err := c.Fetch(ctx, provider, limit)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = 3 * time.Second
	}
	if interval <= 0 {
		interval = ttl
	}
	go notify(items, ttl, interval)
	return nil
}

func notify(items []Item, ttl, interval time.Duration) {
	for _, item := range items {
		msg := item.Title
		if item.Source != "" {
			msg = item.Source + " | " + msg
		}
		dwmblocknotify.PUT(msg, ttl)
		time.Sleep(interval)
	}
}

type rssProvider struct {
	name     string
	source   string
	url      string
	client   *http.Client
	cacheTTL time.Duration
	mu       sync.Mutex
	cache    []Item
	cachedAt time.Time
}

type rssFeed struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

const (
	topNewsURL     = "https://news.google.com/rss/topics/CAAqKggKIiRDQkFTRlFvSUwyMHZNRGx6TVdZU0JYcG9MVU5PR2dKRFRpZ0FQAQ?hl=zh-CN&gl=CN&ceid=CN:zh-Hans"
	financeNewsURL = "https://news.google.com/rss/search?q=%E8%B4%A2%E7%BB%8F&hl=zh-CN&gl=CN&ceid=CN:zh-Hans"
)

func (p *rssProvider) Name() string { return p.name }

func (p *rssProvider) Fetch(ctx context.Context, limit int) ([]Item, error) {
	p.mu.Lock()
	if len(p.cache) > 0 && time.Since(p.cachedAt) < p.cacheTTL {
		items := cloneItems(p.cache)
		p.mu.Unlock()
		return limitItems(items, limit), nil
	}
	p.mu.Unlock()

	items, err := p.fetch(ctx, limit)
	if err != nil {
		p.mu.Lock()
		cached := cloneItems(p.cache)
		p.mu.Unlock()
		if len(cached) > 0 {
			return limitItems(cached, limit), nil
		}
		return nil, err
	}

	p.mu.Lock()
	p.cache = cloneItems(items)
	p.cachedAt = time.Now()
	p.mu.Unlock()
	return limitItems(items, limit), nil
}

func (p *rssProvider) fetch(ctx context.Context, limit int) ([]Item, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s status %d", p.name, resp.StatusCode)
	}
	var feed rssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}
	items := make([]Item, 0, min(limit, len(feed.Channel.Items)))
	for _, raw := range feed.Channel.Items {
		title := strings.TrimSpace(raw.Title)
		if title == "" {
			continue
		}
		items = append(items, Item{Title: title, Source: p.source, Link: raw.Link})
		if len(items) >= limit {
			break
		}
	}
	return items, nil
}

func cloneItems(items []Item) []Item {
	out := make([]Item, len(items))
	copy(out, items)
	return out
}

func limitItems(items []Item, limit int) []Item {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}
