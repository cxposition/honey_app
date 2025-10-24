package grpc_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_server/internal/rpc/node_rpc"
	"io"
)

var CmdRequestChan = make(chan *node_rpc.CmdRequest, 0)

func (NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				logrus.Infof("节点断开")
				return
			}
			if err != nil {
				logrus.Errorf("节点出错: %s", err)
				return
			}
			// 节点往管理发的,命令的执行结果
			fmt.Println("命令结果", response)
		}
	}()

	for request := range CmdRequestChan {
		err := stream.Send(request)
		if err != nil {
			logrus.Errorf("数据发送失败: %s", err)
			continue
		}
	}

	return nil
}
