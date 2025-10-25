package main

import (
	"honey_server/internal/core"
	"honey_server/internal/flags"
	"honey_server/internal/global"
	"honey_server/internal/routers"
	"honey_server/internal/service/grpc_service"
)

func main() {
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()
	core.InitIPDB()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	go grpc_service.Run()
	flags.Run()
	routers.Run()
}
