package main

import (
	"fmt"
	"github.com/j-keck/arping"
	"net"
)

func main() {
	//mac, t, err := arping.Ping(net.ParseIP("111.111.111.133"))
	//fmt.Println(mac, t, err)

	name, duration, err := arping.PingOverIfaceByName(net.ParseIP("111.111.111.1"), "ens37")
	fmt.Println(name, duration, err)
}
