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

	// 2. 初始化 gRPC 客户端
	global.GrpcClient = core.GetGrpcClient()

	// 3. 注册节点
	if err := command.Register(); err != nil {
		logrus.Fatalf("节点注册失败：%v", err)
	}

	// 4. 启动命令监听与任务调度
	go command.RunCommandLoop()
	cron_service.Run()

	// 5. 阻塞主线程
	select {}
}
