package main

import (
	"fmt"
	"honey_server/internal/utils/ip"
)

// 11111111 11111111 10000000 00000000
// 255.255.128.0
func main() {
	_range, err := ip.GetUsableIPRange("192.168.100.6/17")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(_range, err)
}
