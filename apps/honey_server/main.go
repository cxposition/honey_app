package main

import (
	"honey_server/internal/core"
	"honey_server/internal/flags"
	"honey_server/internal/global"
	"honey_server/internal/routers"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/service/mq_service"
)

func main() {
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.Log = core.GetLogger()
	core.InitIPDB()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	global.Queue = core.InitMQ()
	mq_service.RegisterExchange()

	go grpc_service.Run()
	flags.Run()
	routers.Run()
}
