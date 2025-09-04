package main

import (
	"fmt"
	"honey_app/apps/honey_server/internal/core"
)

func main() {
	core.InitIPDB()
	addr := core.GetIpAddr("111.111.111.133")
	fmt.Println(addr)
}
