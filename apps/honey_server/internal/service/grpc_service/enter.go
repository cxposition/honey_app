package grpc_service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"honey_server/internal/global"
	"honey_server/internal/rpc/node_rpc"
	"net"
)

type NodeService struct {
	node_rpc.UnimplementedNodeServiceServer
}

func (NodeService) Register(ctx context.Context, request *node_rpc.RegisterRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	fmt.Println("节点注册:", request)
	return
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
