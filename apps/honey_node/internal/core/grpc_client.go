package core

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"time"
)

var grpcConn *grpc.ClientConn

func GetGrpcClient() node_rpc.NodeServiceClient {
	target := global.Config.System.GrpcManageAddr

	// 如果之前有连接，先关闭它
	if grpcConn != nil {
		_ = grpcConn.Close()
		grpcConn = nil
	}

	var err error
	grpcConn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),                // 阻塞直到连接建立
		grpc.WithTimeout(5*time.Second), // 拨号超时
	)
	if err != nil {
		logrus.Errorf("❌ 连接 gRPC 管理端失败: %v", err)
		return nil
	}

	client := node_rpc.NewNodeServiceClient(grpcConn)
	logrus.Infof("✅ gRPC 客户端连接成功: %s", target)
	return client
}
