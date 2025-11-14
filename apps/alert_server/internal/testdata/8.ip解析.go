package main

import (
	"fmt"
	"net"
	"strings"
)

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

func main() {
	input := "192.168.200.2-192.168.200.200,192.168.200.240"
	ips, err := ParseIPList(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("解析结果数量:", len(ips))
	fmt.Println("前几个:", ips)
	//fmt.Println("最后几个:", ips[len(ips)-3:])
}
