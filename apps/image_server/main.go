package main

import (
	"image_server/internal/core"
	"image_server/internal/flags"
	"image_server/internal/global"
	"image_server/internal/routers"
	"image_server/internal/service/vs_net_service"
)

func main() {
	core.InitIPDB()
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()
	core.SetLogDefault()
	global.DB = core.GetDB()
	global.Redis = core.GetRedis()
	global.DockerClient = core.InitDocker()
	//cron_service.Run()
	vs_net_service.Run()
	flags.Run()
	routers.Run()
}
