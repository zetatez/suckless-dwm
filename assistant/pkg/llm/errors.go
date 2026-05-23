package llm

import "errors"

var (
	ErrNotImplemented = errors.New("llm: not implemented")
	ErrTimeout        = errors.New("llm: request timeout")
	ErrMaxRetries     = errors.New("llm: max retries exceeded")
	ErrInvalidConfig  = errors.New("llm: invalid config")
	ErrUnsupported    = errors.New("llm: unsupported")
	ErrStreamNotReady = errors.New("llm: stream not ready")
)

type ProviderError struct {
	Code    string
	Message string
	Raw     any
}

func (e *ProviderError) Error() string {
	if e.Code != "" {
		return "llm: provider error (" + e.Code + "): " + e.Message
	}
	return "llm: provider error: " + e.Message
}

func IsErrRetryable(err error) bool {
	if e, ok := err.(*ProviderError); ok {
		switch e.Code {
		case "rate_limited", "too_many_requests", "service_unavailable", "server_error":
			return true
		}
	}
	if e, ok := err.(*HTTPError); ok && e.Code >= 500 {
		return true
	}
	return false
}
