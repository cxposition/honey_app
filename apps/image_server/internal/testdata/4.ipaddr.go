package main

import (
	"fmt"
	"image_server/internal/core"
)

func main() {
	core.InitIPDB()
	addr := core.GetIpAddr("111.111.111.133")
	fmt.Println(addr)
}
