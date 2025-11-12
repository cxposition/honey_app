package info

import (
	"fmt"
	"net"
)

// GetNetworkInterfaces 获取网卡信息，key 是网卡名称，value 是 IPv4 地址切片
func GetNetworkInterfaces() (map[string][]string, error) {
	result := make(map[string][]string)

	// 获取所有网卡接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网卡信息失败: %v", err)
	}

	for _, iface := range interfaces {
		// 跳过被禁用或回环网卡
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
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

			// 只取 IPv4 地址
			if ip != nil && ip.To4() != nil {
				result[iface.Name] = append(result[iface.Name], ip.String())
			}
		}
	}

	return result, nil
}
