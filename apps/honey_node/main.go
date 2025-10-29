package main

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/service/command"
	"honey_node/internal/service/cron_service"
)

func main() {
	// 1. 初始化配置与日志
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.Log = core.GetLogger()
	logrus.Infof("启动成功，版本：%s, 提交：%s", global.Version, global.Commit)

	// 2. 启动节点注册与命令循环（包含 gRPC 自动重连）
	go command.RunCommandLoop()

	// 3. 启动定时任务（资源上报）
	go cron_service.Run()

	// 4. 阻塞主线程，防止退出
	select {}
}
