package main

import (
	"honey_app/apps/honey_server/core"
	"honey_app/apps/honey_server/flags"
	"honey_app/apps/honey_server/global"
)

func main() {
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	flags.Run()
}
