package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BaseClient struct {
	baseURL     string
	httpClient  *http.Client
	maxRetries  int
	backoffBase int
	backoffType BackoffType
}

func NewBaseClient(baseURL string, cfg Config) *BaseClient {
	return &BaseClient{
		baseURL:     baseURL,
		httpClient:  cfg.GetHTTPClient(),
		maxRetries:  cfg.GetMaxRetries(),
		backoffBase: cfg.GetBackoffBase(),
		backoffType: cfg.GetBackoffType(),
	}
}

func (c *BaseClient) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *BaseClient) Do(ctx context.Context, method, path string, body any, headers map[string]string) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.calculateBackoff(attempt)
			time.Sleep(backoff)
		}

		var reqBody []byte
		var err error
		if body != nil {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
		if err != nil {
			return nil, err
		}

		if len(reqBody) > 0 {
			req.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			resp.Body.Close()
			lastErr = &HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("server error: %d", resp.StatusCode)}
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

func (c *BaseClient) calculateBackoff(attempt int) time.Duration {
	if c.backoffType == BackoffExponential {
		return time.Duration(c.backoffBase*c.backoffBase*attempt) * time.Second
	}
	return time.Duration(c.backoffBase*attempt) * time.Second
}

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}
