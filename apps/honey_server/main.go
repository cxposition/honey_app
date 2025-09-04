package main

import (
	"honey_app/apps/honey_server/internal/core"
	"honey_app/apps/honey_server/internal/flags"
	"honey_app/apps/honey_server/internal/global"
	"honey_app/apps/honey_server/internal/routers"
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
