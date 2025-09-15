package core

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func InitDocker() *client.Client {
	// 创建docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Fatalf("创建docker客户端失败: %v", err)
	}
	logrus.Infof("创建docker客户端成功")
	return cli
}
