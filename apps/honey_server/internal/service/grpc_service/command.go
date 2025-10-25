package grpc_service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"honey_server/internal/rpc/node_rpc"
	"io"
)

type Command struct {
	ReqChan chan *node_rpc.CmdRequest
	ResChan chan *node_rpc.CmdResponse
	Server  node_rpc.NodeService_CommandServer
}

var NodeCommandMap = map[string]*Command{}

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
	NodeCommandMap[nodeID] = &Command{
		ReqChan: make(chan *node_rpc.CmdRequest),
		ResChan: make(chan *node_rpc.CmdResponse),
		Server:  stream,
	}

	go func() {
		for request := range NodeCommandMap[nodeID].ReqChan {
			err := NodeCommandMap[nodeID].Server.Send(request)
			if err != nil {
				logrus.Errorf("数据发送失败: %s", err)
				continue
			}
		}
	}()

	for {
		response, err := NodeCommandMap[nodeID].Server.Recv()
		if err == io.EOF {
			logrus.Infof("节点断开")
			break
		}
		if err != nil {
			logrus.Errorf("节点出错: %s", err)
			break
		}
		// 节点往管理发的,命令的执行结果
		fmt.Println("命令结果", response)
		NodeCommandMap[nodeID].ResChan <- response
	}

	delete(NodeCommandMap, nodeID)
	return nil
}
