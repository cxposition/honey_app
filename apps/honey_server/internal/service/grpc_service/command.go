package grpc_service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"honey_server/internal/rpc/node_rpc"
	"io"
)

var CmdRequestChan = make(chan *node_rpc.CmdRequest, 0)
var CmdResponseChan = make(chan *node_rpc.CmdResponse, 0)

var streamMap = map[string]node_rpc.NodeService_CommandServer{}

func (NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	nodeIDList := md.Get("NodeID")
	if len(nodeIDList) == 0 {
		return errors.New("请在metadata中传入节点id")
	}
	nodeID := nodeIDList[0]
	streamMap[nodeID] = stream
	go func() {
		for {
			response, err := streamMap[nodeID].Recv()
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
		fmt.Println("stream2:", stream)
		err := streamMap[nodeID].Send(request)
		if err != nil {
			logrus.Errorf("数据发送失败: %s", err)
			continue
		}
	}

	return nil
}
