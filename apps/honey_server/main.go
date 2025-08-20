package main

import (
	"fmt"
	"honey_app/apps/honey_server/core"
	"honey_app/apps/honey_server/flags"
	"honey_app/apps/honey_server/global"
)

func main() {
	global.Config = core.ReadConfig()
	fmt.Println(global.Config)
	global.DB = core.InitDB()
	flags.Run()
}
