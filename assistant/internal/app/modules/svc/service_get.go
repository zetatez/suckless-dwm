package svc

import (
	"fmt"
	"net"
	"time"

	"assistant/internal/bootstrap/psl"
)

func (s *Service) GetIP() ([]string, error) {
	iface := psl.GetConfig().App.Interface
	netIface, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, fmt.Errorf("get interface %s: %w", iface, err)
	}
	addrs, err := netIface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("get addresses: %w", err)
	}
	var ips []string
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if !ip.IsLoopback() && ip.To4() != nil {
			ips = append(ips, ip.String())
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPv4 addresses found on %s", iface)
	}
	s.copyToClipboardWithNotify(ips[0], fmt.Sprintf("get success: %s", ips[0]))
	return ips, nil
}

func (s *Service) GetDatetime() map[string]string {
	now := time.Now()
	result := map[string]string{
		"datetime":   now.Format(time.DateTime),
		"unix":       fmt.Sprintf("%d", now.Unix()),
		"unix_milli": fmt.Sprintf("%d", now.UnixMilli()),
		"date":       now.Format("2006-01-02"),
		"time":       now.Format("15:04:05"),
	}
	s.copyToClipboardWithNotify(result["datetime"], fmt.Sprintf("get success: %s", result["datetime"]))
	return result
}

func (s *Service) GetUnixSec() string {
	now := fmt.Sprintf("%d", time.Now().Unix())
	s.copyToClipboardWithNotify(now, fmt.Sprintf("get success: %s", now))
	return now
}
