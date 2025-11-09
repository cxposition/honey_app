package grpc_service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
)

func (NodeService) StatusCreateIP(ctx context.Context, in *node_rpc.StatusCreateRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	var honeyIPModel models.HoneyIpModel
	err1 := global.DB.Take(&honeyIPModel, in.HoneyIPID).Error
	if err1 != nil {
		return nil, fmt.Errorf("诱捕ip不存在 %d", in.HoneyIPID)
	}

	var status int8 = 2
	if in.ErrMsg != "" {
		status = 3
		logrus.Errorf("创建ip失败 %s", in.ErrMsg)
	}

	global.DB.Model(&honeyIPModel).Updates(models.HoneyIpModel{
		Mac:     in.Mac,
		Network: in.Network,
		Status:  status,
	})
	return
}
