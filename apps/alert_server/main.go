package main

import (
	"alert_server/internal/core"
	"alert_server/internal/flags"
	"alert_server/internal/global"
	"alert_server/internal/routers"
	"alert_server/internal/service/mq_service"
)

func main() {
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.Log = core.GetLogger()
	core.InitIPDB()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	global.MQConn = core.InitMQ()
	global.ES = core.InitES()
	go mq_service.Run()

	flags.Run()
	routers.Run()
}
