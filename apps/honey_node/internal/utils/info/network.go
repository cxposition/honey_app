package info

import (
	"net"
	"strings"
)

type NetworkInfo struct {
	Network string
	Ip      string
	Mask    int
	Net     string
}

func GetNetworkList(filterNetworkList []string) (list []NetworkInfo, err error) {
	faces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, face := range faces {
		faceName := face.Name
		if faceName == "lo" {
			continue
		}

		var isFilter bool
		// 过滤掉诱捕ip的网卡
		for _, s := range filterNetworkList {
			if strings.HasPrefix(faceName, s) {
				isFilter = true
				break
			}
		}

		if isFilter {
			continue
		}

		addrs, err := face.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, _net, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			//if ip.To4() == nil {
			//	continue
			//}

			mask, _ := _net.Mask.Size()
			list = append(list, NetworkInfo{
				Network: faceName,
				Ip:      ip.String(),
				Mask:    mask,
				Net:     _net.String(),
			})
		}
	}

	return
}
