package main

import (
	"fmt"
	"honey_node/internal/utils/ip"
)

func main() {
	ifaceName := "ens37"
	ipAddr, mac, err := ip.GetNetworkInfo(ifaceName)
	fmt.Println(ipAddr, mac, err)
}
