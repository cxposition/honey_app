package port_service

import (
	"context"
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"io"
	"net"
)

func Tunnel(localAddr string, targetAddr string) error {
	// 创建本地监听
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		logrus.Errorf("创建本地监听失败: %v", err)
		return err
	}
	defer listener.Close()

	logrus.Infof("本地监听启动，地址: %s", localAddr)
	logrus.Infof("目标地址: %s", targetAddr)

	// 接受客户端连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("接受客户端连接失败: %v", err)
			continue
		}

		// 为每个连接创建一个goroutine处理
		go func() {
			err := handleConnection(global.GrpcClient, clientConn, targetAddr)
			if err != nil {
				logrus.Errorf("处理连接失败: %v", err)
			}
		}()
	}
}

func handleConnection(client node_rpc.NodeServiceClient, localConn net.Conn, targetAddr string) error {
	defer localConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建双向流
	stream, err := client.Tunnel(ctx)
	if err != nil {
		logrus.Errorf("创建隧道失败: %v", err)
		return err
	}

	// 发送初始消息，包含目标地址
	if err := stream.Send(&node_rpc.TunnelData{
		Chunk:   []byte{},
		Address: targetAddr,
	}); err != nil {
		logrus.Errorf("发送初始请求失败: %v", err)
		return err
	}

	// 创建用于接收gRPC服务器数据并发往本地连接的goroutine
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logrus.Errorf("接收gRPC服务器数据失败: %v", err)
				break
			}

			// 将数据写入本地连接
			_, err = localConn.Write(resp.Chunk)
			if err != nil {
				logrus.Errorf("写入本地连接失败: %v", err)
				break
			}
		}
		cancel()
	}()

	// 从本地连接读取数据并发送到gRPC服务器
	buffer := make([]byte, 4096)
	for {
		n, err := localConn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logrus.Infof("本地连接已关闭")
			} else {
				logrus.Errorf("从本地连接读取失败: %v", err)
			}
			break
		}

		// 发送数据到gRPC服务器
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n],
			Address: targetAddr,
		})
		if err != nil {
			logrus.Errorf("发送数据到gRPC服务器失败: %v", err)
			break
		}
	}

	// 关闭流
	_ = stream.CloseSend()
	return nil
}
