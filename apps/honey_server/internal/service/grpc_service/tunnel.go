package grpc_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_server/internal/rpc/node_rpc"
	"io"
	"log"
	"net"
)

// Tunnel 实现双向流 RPC
func (s *NodeService) Tunnel(stream node_rpc.NodeService_TunnelServer) error {
	// 接收客户端的第一个消息，获取目标地址
	//logrus.Infof("Tunnel server 执行 stream %p", stream)	// 输出stream的内存地址, stream内存地址每次都不一样
	req, err := stream.Recv()
	if err != nil {
		logrus.Errorf("接收初始请求失败: %v", err)
		return fmt.Errorf("接收初始请求失败: %v", err)
	}

	// 连接到目标地址, 这是创建一个tcp客户端
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(stream.Context(), "tcp", req.Address)
	if err != nil {
		logrus.Errorf("连接目标地址失败: %v", err)
		return fmt.Errorf("连接目标地址失败: %v", err)
	}
	defer conn.Close()
	logrus.Infof("已连接至目标地址: %s", req.Address)

	// 创建用于接收客户端数据并发往目标地址的goroutine
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("接收客户端数据失败: %v", err)
				return
			}

			// 将数据写入目标连接
			_, err = conn.Write(req.Chunk)
			if err != nil {
				log.Printf("写入目标连接失败: %v", err)
				return
			}
		}
	}()

	// 从目标地址读取数据并发送给客户端
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("目标连接已关闭")
			} else {
				log.Printf("从目标连接读取失败: %v", err)
			}
			return nil
		}

		// 发送数据到客户端
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n],
			Address: req.Address,
		})
		if err != nil {
			log.Printf("发送数据到客户端失败: %v", err)
			return err
		}
	}
}
