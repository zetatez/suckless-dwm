package svc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *Service) ConvertDatetime(from, to string) (string, error) {
	clip, err := s.readClipboard()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}
	input := strings.TrimSpace(clip)
	if input == "" {
		return "", fmt.Errorf("clipboard is empty")
	}

	var t time.Time

	switch from {
	case "datetime":
		t, err = time.Parse(time.DateTime, strings.TrimSpace(input))
	case "unix":
		sec, e := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if e != nil {
			return "", fmt.Errorf("invalid unix timestamp: %w", e)
		}
		t = time.Unix(sec, 0)
	case "unix_milli":
		ms, e := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if e != nil {
			return "", fmt.Errorf("invalid unix milli timestamp: %w", e)
		}
		t = time.UnixMilli(ms)
	default:
		return "", fmt.Errorf("unsupported from format: %s, available: datetime, unix, unix_milli", from)
	}
	if err != nil {
		return "", fmt.Errorf("parse input failed: %w", err)
	}

	var result string
	switch to {
	case "datetime":
		result = t.Format(time.DateTime)
	case "unix":
		result = fmt.Sprintf("%d", t.Unix())
	case "unix_milli":
		result = fmt.Sprintf("%d", t.UnixMilli())
	default:
		return "", fmt.Errorf("unsupported to format: %s, available: datetime, unix, unix_milli", to)
	}
	_, err = s.pushClipboard(result, fmt.Sprintf("transfer success: %s", result))
	return result, err
}
