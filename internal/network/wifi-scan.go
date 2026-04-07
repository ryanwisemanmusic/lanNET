package network

import (
	"fmt"
	"net"
	"strings"
)

func GetLocalIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if (iface.Flags&net.FlagUp) == 0 ||
			(iface.Flags&net.FlagLoopback) != 0 {
			continue
		}

		if !looksLikeWifiOrLAN(iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}

	}
	return "", fmt.Errorf("no active Wi-Fi/LAN IPv4 address found")
}

func looksLikeWifiOrLAN(name string) bool {
	name = strings.ToLower(name)

	return strings.HasPrefix(name, "en") ||
		strings.HasPrefix(name, "eth") ||
		strings.HasPrefix(name, "wi-fi") ||
		strings.HasPrefix(name, "wifi")
}
