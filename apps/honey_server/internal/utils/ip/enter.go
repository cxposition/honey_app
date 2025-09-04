package ip

import "net"

func HasLocalIPAddr(_ip string) bool {
	ip := net.ParseIP(_ip)
	if ip.IsLoopback() {
		return true
	}
	if ip.IsLinkLocalUnicast() {
		return true
	}
	return false
}
