package grpc_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
)

func (NodeService) Register(ctx context.Context, request *node_rpc.RegisterRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	// 节点不存在，需要创建
	uid := request.NodeUid
	var model models.NodeModel
	err1 := global.DB.Take(&model, "uid = ?", uid).Error
	if err1 != nil {
		// 创建节点
		model = models.NodeModel{
			Title:  request.SystemInfo.HostName,
			Uid:    uid,
			IP:     request.Ip,
			Status: 1,
			SystemInfo: models.NodeSystemInfo{
				NodeVersion: request.Version,
				NodeCommit:  request.Commit,
				//HostName:
				//DistributionVersion:
				//CoreVersion:
				//SystemType:
				//StartTime:
			},
		}
		err1 = global.DB.Create(&model).Error
		if err1 != nil {
			logrus.Errorf("节点创建失败 %s", err1)
			return nil, errors.New("节点创建失败")
		}
	}

	if model.Status != 1 {
		// 改状态
		global.DB.Model(&model).Update("status", 1)
	}

	// 节点存在
	fmt.Println("节点注册:", request)
	return
}
