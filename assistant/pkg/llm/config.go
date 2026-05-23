package llm

import (
	"net/http"
	"time"
)

type BackoffType int

const (
	BackoffLinear BackoffType = iota
	BackoffExponential
)

const (
	DefaultTimeout     = 30
	DefaultMaxRetries  = 3
	DefaultBackoffBase = 1
	DefaultBackoffType = BackoffExponential
	DefaultMaxTokens   = 4096
	DefaultTemperature = 0.7
)

type Config struct {
	APIKey  string
	BaseURL string
	Model   string
	Extra   map[string]string

	Timeout     int
	MaxRetries  int
	BackoffBase int
	BackoffType BackoffType
	HTTPClient  *http.Client

	MaxTokens   int
	Temperature float32
	TopP        float32
}

func (c Config) GetTimeout() time.Duration {
	if c.Timeout <= 0 {
		return DefaultTimeout * time.Second
	}
	return time.Duration(c.Timeout) * time.Second
}

func (c Config) GetMaxRetries() int {
	if c.MaxRetries <= 0 {
		return DefaultMaxRetries
	}
	return c.MaxRetries
}

func (c Config) GetBackoffBase() int {
	if c.BackoffBase <= 0 {
		return DefaultBackoffBase
	}
	return c.BackoffBase
}

func (c Config) GetBackoffType() BackoffType {
	if c.BackoffType == 0 {
		return DefaultBackoffType
	}
	return c.BackoffType
}

func (c Config) GetMaxTokens() int {
	if c.MaxTokens <= 0 {
		return DefaultMaxTokens
	}
	return c.MaxTokens
}

func (c Config) GetTemperature() float32 {
	if c.Temperature == 0 {
		return DefaultTemperature
	}
	return c.Temperature
}

func (c Config) GetHTTPClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{
		Timeout: c.GetTimeout(),
	}
}
