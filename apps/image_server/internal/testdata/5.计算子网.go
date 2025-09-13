package main

import (
	"fmt"
	"net"
)

func main() {
	x := "10.2.2.0/16"
	ip, ipNet, err := net.ParseCIDR(x)
	fmt.Println("ip:", ip, "ipNet:", ipNet, "err:", err)
	ip4 := ip.To4()
	fmt.Println(ip4[0], ip4[1], ip4[2], ip4[3])
}
