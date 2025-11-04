package ip

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

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

func ParseIPList(ipStr string) ([]string, error) {
	var result []string
	parts := strings.Split(ipStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 判断是否是 IP 段
		if strings.Contains(part, "-") {
			ipRange := strings.Split(part, "-")
			if len(ipRange) != 2 {
				return nil, fmt.Errorf("invalid ip range: %s", part)
			}
			startIP := net.ParseIP(strings.TrimSpace(ipRange[0])).To4()
			endIP := net.ParseIP(strings.TrimSpace(ipRange[1])).To4()
			if startIP == nil || endIP == nil {
				return nil, fmt.Errorf("invalid ip: %s", part)
			}

			// 展开 IP 段
			for ip := cloneIP(startIP); !ipAfter(ip, endIP); incIP(ip) {
				result = append(result, ip.String())
			}
		} else {
			ip := net.ParseIP(part).To4()
			if ip == nil {
				return nil, fmt.Errorf("invalid ip: %s", part)
			}
			result = append(result, ip.String())
		}
	}

	return result, nil
}

// cloneIP 拷贝一个 IP
func cloneIP(ip net.IP) net.IP {
	ipv4 := make(net.IP, len(ip))
	copy(ipv4, ip)
	return ipv4
}

// incIP 将 IP +1
func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

// ipAfter 判断 ip1 是否在 ip2 之后
func ipAfter(ip1, ip2 net.IP) bool {
	for i := 0; i < len(ip1); i++ {
		if ip1[i] > ip2[i] {
			return true
		}
		if ip1[i] < ip2[i] {
			return false
		}
	}
	return false
}

// IncrementIP 让 IP 地址 +1，例如 192.168.1.1 -> 192.168.1.2
func IncrementIP(ip net.IP) net.IP {
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)

	for i := len(newIP) - 1; i >= 0; i-- {
		newIP[i]++
		if newIP[i] != 0 {
			break
		}
	}
	return newIP
}

// DecrementIP 让 IP 地址 -1，例如 192.168.1.255 -> 192.168.1.254
func DecrementIP(ip net.IP) net.IP {
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)

	for i := len(newIP) - 1; i >= 0; i-- {
		if newIP[i] == 0 {
			newIP[i] = 255
		} else {
			newIP[i]--
			break
		}
	}
	return newIP
}

// FormatIPRange 格式化IP范围，并自动跳过 .1 和 .2
func FormatIPRange(startIP, endIP net.IP) string {
	startStr := startIP.String()
	endStr := endIP.String()

	// 自动检测是否是常见的私有网段（192.168.x.x）
	if strings.HasPrefix(startStr, "192.168.") {
		parts := strings.Split(startStr, ".")
		if len(parts) == 4 {
			// 把 .1 和 .2 跳过，从 .3 开始
			startStr = fmt.Sprintf("%s.%s.%s.3", parts[0], parts[1], parts[2])
		}
	}

	return fmt.Sprintf("%s-%s", startStr, endStr)
}

// BroadcastIP 计算子网的广播地址（便于 DecrementIP 使用）
// 比如 192.168.1.0/24 -> 192.168.1.255
func BroadcastIP(ipNet *net.IPNet) net.IP {
	ip := ipNet.IP.To4()
	mask := ipNet.Mask
	if ip == nil {
		return nil
	}

	broadcast := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast
}

// ParseIPRange 把 "192.168.1.1-192.168.1.254" 转成 IP 列表
func ParseIPRange(ipRange string) ([]string, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的 IP 范围: %s", ipRange)
	}

	start := net.ParseIP(strings.TrimSpace(parts[0])).To4()
	end := net.ParseIP(strings.TrimSpace(parts[1])).To4()
	if start == nil || end == nil {
		return nil, fmt.Errorf("无效的 IPv4 地址")
	}

	startInt := binary.BigEndian.Uint32(start)
	endInt := binary.BigEndian.Uint32(end)

	if endInt < startInt {
		return nil, fmt.Errorf("结束 IP 小于起始 IP")
	}

	var ips []string
	for i := startInt; i <= endInt; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}
	return ips, nil
}
