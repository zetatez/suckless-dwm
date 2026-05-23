package middleware

import (
	"net/http"
	"sync"
	"time"

	"assistant/pkg/response"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests    map[string][]time.Time
	mu          sync.RWMutex
	maxRequests int
	window      time.Duration
	stopCh      chan struct{}
	stopped     bool
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: maxRequests,
		window:      window,
		stopCh:      make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.stopped {
		return
	}
	rl.stopped = true
	close(rl.stopCh)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-rl.stopCh:
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, times := range rl.requests {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) <= rl.window {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(rl.requests, key)
				} else {
					rl.requests[key] = valid
				}
			}
			rl.mu.Unlock()
		}
	}
}

func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[key]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) <= rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.maxRequests {
		rl.requests[key] = valid
		return false
	}

	valid = append(valid, now)
	rl.requests[key] = valid
	return true
}

var globalLimiter *RateLimiter

func RateLimit(maxRequests int, windowSecs int) gin.HandlerFunc {
	limiter := NewRateLimiter(maxRequests, time.Duration(windowSecs)*time.Second)
	globalLimiter = limiter
	return limiter.Handle()
}

func StopLimiter() {
	if globalLimiter != nil {
		globalLimiter.Stop()
	}
}

func (rl *RateLimiter) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if !rl.isAllowed(key) {
			response.Err(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
