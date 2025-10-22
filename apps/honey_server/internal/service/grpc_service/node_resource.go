package grpc_service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
)

func (NodeService) NodeResource(ctx context.Context, request *node_rpc.NodeResourceRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	uid := request.NodeUid
	var model models.NodeModel
	err1 := global.DB.Take(&model, "uid = ?", uid).Error
	if err1 != nil {
		return nil, errors.New("节点不存在")
	}
	newModel := models.NodeModel{
		Resource: models.NodeResource{
			CpuCount:              int(request.ResourceInfo.CpuCount),
			CpuUseRate:            float64(request.ResourceInfo.CpuUseRate),
			MemTotal:              request.ResourceInfo.MemTotal,
			MemUseRate:            float64(request.ResourceInfo.MemUseRate),
			DiskTotal:             request.ResourceInfo.DiskTotal,
			DiskUseRate:           float64(request.ResourceInfo.DiskUseRate),
			NodePath:              request.ResourceInfo.NodePath,
			NodeResourceOccupancy: request.ResourceInfo.NodeResourceOccupancy,
		},
	}

	err1 = global.DB.Model(&model).Updates(newModel).Error
	if err1 != nil {
		logrus.Errorf("节点资源状态更新失败 %s", err)
		return nil, errors.New("节点资源状态更新失败")
	}
	return
}
