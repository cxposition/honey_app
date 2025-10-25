package cron_service

import (
	"context"
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
)

func Resource() {
	if global.GrpcClient == nil {
		logrus.Errorf("管理未连接，放弃上报")
		return
	}

	resourceInfo, err := info.GetSystemResource()
	if err != nil {
		logrus.Fatalln(err)
	}

	_, err = global.GrpcClient.NodeResource(context.Background(), &node_rpc.NodeResourceRequest{
		NodeUid: global.Config.System.Uid,
		ResourceInfo: &node_rpc.ResourceMessage{
			CpuCount:              resourceInfo.CPUCount,
			CpuUseRate:            float32(resourceInfo.CPUUseRate),
			MemTotal:              resourceInfo.MemTotal,
			MemUseRate:            float32(resourceInfo.MemUseRate),
			DiskTotal:             resourceInfo.DiskTotal,
			DiskUseRate:           float32(resourceInfo.DiskUseRate),
			NodePath:              resourceInfo.NodePath,
			NodeResourceOccupancy: resourceInfo.NodeResourceOccupancy,
		},
	})
	if err != nil {
		logrus.Errorf("上报资源信息失败：%s", err)
	}
	logrus.Infof("上报资源信息成功")
}
