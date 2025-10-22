package core

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/rpc"
	"honey_node/internal/rpc/node_rpc"
)

func GetGrpcClient() node_rpc.NodeServiceClient {
	addr := global.Config.System.GrpcManageAddr
	conn := rpc.GetConn(addr)
	// 初始化客户端
	client := node_rpc.NewNodeServiceClient(conn)
	logrus.Infof("初始化grpc客户端成功 %s", addr)
	return client
}
