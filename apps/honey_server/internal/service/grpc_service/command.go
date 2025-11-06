package grpc_service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"honey_server/internal/rpc/node_rpc"
	"io"
	"sync"
	"time"
)

// Command 表示一个节点的 gRPC 双向流连接
type Command struct {
	ReqChan chan *node_rpc.CmdRequest
	ResMap  sync.Map // taskID -> chan *node_rpc.CmdResponse
	Server  node_rpc.NodeService_CommandServer
	CloseCh chan struct{}
	NodeID  string
}

// NodeCommandMap 存放所有在线节点
var NodeCommandMap sync.Map // nodeID -> *Command
var mapMutex sync.RWMutex

func (NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("missing metadata")
	}

	nodeIDList := md.Get("NodeID")
	if len(nodeIDList) == 0 {
		return errors.New("请在metadata中传入节点id")
	}
	nodeID := nodeIDList[0]

	cmd := &Command{
		ReqChan: make(chan *node_rpc.CmdRequest, 10),
		Server:  stream,
		CloseCh: make(chan struct{}),
		NodeID:  nodeID,
	}

	NodeCommandMap.Store(nodeID, cmd)
	logrus.Infof("节点 %s 已连接", nodeID)

	// 异步发送循环
	go func(c *Command) {
		for {
			select {
			case <-c.CloseCh:
				return
			case req := <-c.ReqChan:
				if err := c.Server.Send(req); err != nil {
					logrus.Errorf("发送到节点 %s 失败: %v", c.NodeID, err)
					return
				}
			}
		}
	}(cmd)

	// 接收循环
	for {
		resp, err := stream.Recv()
		logrus.Infof("节点 %s 接收到命令: %+v", nodeID, resp)
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("节点 %s 接收错误: %v", nodeID, err)
			break
		}

		// 匹配 taskID 返回结果
		if chVal, ok := cmd.ResMap.Load(resp.TaskID); ok {
			ch := chVal.(chan *node_rpc.CmdResponse)
			select {
			case ch <- resp: // 匹配好结果后将结果发送到对应通道,相当于向respChan中发消息
				logrus.Infof("节点 %s 的 task %s 响应已发送", nodeID, resp.TaskID)
			case <-ctx.Done():
				//default:
				//	logrus.Warnf("节点 %s 的 task %s 响应通道已满", nodeID, resp.TaskID)
			}
		} else {
			logrus.Warnf("节点 %s 收到未知 taskID=%s 的响应", nodeID, resp.TaskID)
		}
	}

	close(cmd.CloseCh)
	NodeCommandMap.Delete(nodeID)
	logrus.Infof("节点 %s 已断开连接", nodeID)
	return nil
}

// SendCommand 发送命令并等待响应（带超时）
func SendCommand(nodeID string, req *node_rpc.CmdRequest, timeout time.Duration) (*node_rpc.CmdResponse, error) {
	val, ok := NodeCommandMap.Load(nodeID)
	if !ok {
		return nil, fmt.Errorf("节点 %s 未连接", nodeID)
	}
	cmd := val.(*Command)

	respChan := make(chan *node_rpc.CmdResponse, 1)
	cmd.ResMap.Store(req.TaskID, respChan)
	defer cmd.ResMap.Delete(req.TaskID)

	select {
	case cmd.ReqChan <- req:
		// 成功发送
		logrus.Infof("节点 %s 的 task %s 已发送", nodeID, req.TaskID)
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("发送命令到节点 %s 超时", nodeID)
	}

	select {
	case resp := <-respChan:
		logrus.Infof("节点 %s 的 task %s 已接收到响应", nodeID, req.TaskID)
		return resp, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("等待节点 %s 响应超时", nodeID)
	}
}

func GetNodeCommand(nodeID string) (*Command, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	cmd, ok := NodeCommandMap.Load(nodeID)
	if !ok {
		return nil, false
	}
	return cmd.(*Command), ok
}
