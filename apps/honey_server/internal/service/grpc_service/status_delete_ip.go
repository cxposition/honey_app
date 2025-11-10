package grpc_service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
)

func (NodeService) StatusDeleteIP(ctx context.Context, in *node_rpc.StatusDeleteIPRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	var honeyIPList []models.HoneyIpModel
	global.DB.Find(&honeyIPList, "id in ?", in.HoneyIPIDList)
	if len(honeyIPList) == 0 {
		return nil, fmt.Errorf("诱捕ip不存在")
	}
	global.DB.Delete(&honeyIPList)
	logrus.Infof("删除诱捕ip成功,ip:%v", in.HoneyIPIDList)
	return
}
