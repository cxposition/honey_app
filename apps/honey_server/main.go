package main

import (
	"honey_app/apps/honey_server/core"
	"honey_app/apps/honey_server/flags"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/routers"
)

func main() {
	core.InitIPDB()
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	flags.Run()
	routers.Run()
}
