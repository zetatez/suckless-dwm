package health

import (
	"time"
)

type HealthStatus struct {
	Status   string           `json:"status"`
	Duration string           `json:"duration"`
	Checks   map[string]Check `json:"checks"`
}

type Check struct {
	Status   string `json:"status"`
	Duration string `json:"duration"`
	Error    string `json:"error,omitempty"`
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

type HealthService struct{}

func (s *HealthService) Health() (*HealthStatus, error) {
	status := &HealthStatus{
		Status: "ok",
		Checks: make(map[string]Check),
	}
	start := time.Now()
	status.Duration = time.Since(start).String()
	return status, nil
}
