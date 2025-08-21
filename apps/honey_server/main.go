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
	flags.Run()
	global.Log.Infof("info日志")
	global.Log.Warnf("warn日志")
	global.Log.Errorf("error日志")
	core.GetDB()
	core.GetDB()
}
