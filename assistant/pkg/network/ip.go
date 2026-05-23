package network

import (
	"net"
)

func getSVRIP(deviceName string) string {
	defaultSVRIP := "127.0.0.1"
	ifs, err := net.Interfaces()
	if err != nil {
		return defaultSVRIP
	}
	for _, iface := range ifs {
		if iface.Name == deviceName && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return defaultSVRIP
}
