package main

import (
	"github.com/sirupsen/logrus"
	"net"
)

func main() {
	//s := "192.168.26.3/24"
	//ip, ipNet, err := net.ParseCIDR(s)
	//fmt.Println(ip, ipNet, err)

	s := "192.168.26.x"
	ip := net.ParseIP(s)
	if ip == nil {
		logrus.Errorf("ip formart is not correct.")
	}
}
