package main

import (
	"fmt"
	"github.com/j-keck/arping"
	"honey_node/internal/core"
	"net"
)

func main() {
	//fmt.Println(core.GetVendorByMAC("00:0c:29:c3:bb:6b"))
	mac, duration, err := arping.PingOverIfaceByName(net.ParseIP("192.168.177.1"), "ens33")
	fmt.Println(mac, duration, err)
	byMAC := core.GetVendorByMAC(mac.String())
	fmt.Println(byMAC)
}
