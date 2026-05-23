package retry

import (
	"context"
	"errors"
	"log"
	"math"
	"time"
)

type Hook struct {
	BeforeRetry func(attempt int, err error)
	AfterRetry  func(attempt int, err error)
	OnSuccess   func(result any)
	OnFail      func(finalErr error)
}

type Config struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	RetryIf      func(error) bool
	Logger       *log.Logger
	Hooks        Hook
}

func DefaultConfig() Config {
	return Config{
		MaxRetries:   3,
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     3 * time.Second,
		Multiplier:   2.0,
		RetryIf:      func(err error) bool { return err != nil }, // Default: retry if err != nill
		Logger:       log.Default(),
	}
}

func Do[T any](ctx context.Context, cfg Config, fn func() (T, error)) (T, error) {
	var zero T
	delay := cfg.InitialDelay

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			cfg.Logger.Println("[retry] context canceled:", ctx.Err())
			return zero, ctx.Err()
		default:
		}

		result, err := fn()
		if err == nil {
			if cfg.Hooks.OnSuccess != nil {
				cfg.Hooks.OnSuccess(result)
			}
			return result, nil
		}

		if cfg.Hooks.BeforeRetry != nil {
			cfg.Hooks.BeforeRetry(attempt, err)
		}

		if cfg.RetryIf != nil && !cfg.RetryIf(err) {
			if cfg.Hooks.OnFail != nil {
				cfg.Hooks.OnFail(err)
			}
			return zero, err
		}

		cfg.Logger.Printf("[retry] attempt=%d failed: %v\n", attempt+1, err)

		if attempt == cfg.MaxRetries {
			if cfg.Hooks.OnFail != nil {
				cfg.Hooks.OnFail(err)
			}
			return zero, err
		}

		if cfg.Hooks.AfterRetry != nil {
			cfg.Hooks.AfterRetry(attempt, err)
		}

		time.Sleep(delay)
		delay = time.Duration(math.Min(float64(cfg.MaxDelay), float64(delay)*cfg.Multiplier))
	}

	return zero, errors.New("unreachable")
}
