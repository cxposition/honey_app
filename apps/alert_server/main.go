package main

import (
	"alert_server/internal/core"
	"alert_server/internal/flags"
	"alert_server/internal/global"
	"alert_server/internal/routers"
)

func main() {
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.Log = core.GetLogger()
	core.InitIPDB()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	//global.MQConn = core.InitMQ()
	global.ES = core.InitES()

	flags.Run()
	routers.Run()
}
