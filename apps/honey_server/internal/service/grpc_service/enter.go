package grpc_service

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"honey_server/internal/global"
	"honey_server/internal/rpc/node_rpc"
	"net"
)

type NodeService struct {
	node_rpc.UnimplementedNodeServiceServer
}

func Run() {
	// 监听端口
	addr := global.Config.System.GrpcAddr
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}

	// 创建一个gRPC服务器实例
	s := grpc.NewServer()
	server := NodeService{}
	// 将server结构体注册为gRPC服务
	node_rpc.RegisterNodeServiceServer(s, &server)
	logrus.Infof("grpc server running %s", addr)
	// 开始处理客户端请求。
	err = s.Serve(listen)
}
