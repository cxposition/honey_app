package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/ip"
	"log"
	"os"
)

func main() {
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()

	addr := global.Config.System.GrpcManageAddr
	// 使用 grpc.Dial 创建一个到指定地址的 gRPC 连接。
	// 此处使用不安全的证书来实现 SSL/TLS 连接
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	defer conn.Close()
	// 初始化客户端
	client := node_rpc.NewNodeServiceClient(conn)
	_ip, mac, err := ip.GetNetworkInfo(global.Config.System.Network)
	if err != nil {
		logrus.Fatalln(err)
	}
	uid := uuid.New().String()
	if global.Config.System.Uid == "" {
		global.Config.System.Uid = uid
	}
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Fatalln(err)
	}
	core.SetConfig()
	result, err := client.Register(context.Background(), &node_rpc.RegisterRequest{
		Ip:      _ip,
		Mac:     mac,
		NodeUid: global.Config.System.Uid,
		Version: global.Version,
		Commit:  global.Commit,
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName:            hostName,
			DistributionVersion: "",
			CoreVersion:         "",
			SystemType:          "",
			StartTime:           "",
		},
	})
	fmt.Println(result, err)
}
