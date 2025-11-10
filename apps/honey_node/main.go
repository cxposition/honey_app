package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/command"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/service/mq_service"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化配置与日志
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.Log = core.GetLogger()
	logrus.Infof("启动成功，版本：%s, 提交：%s", global.Version, global.Commit)

	// 启动节点注册与命令循环（包含 gRPC 自动重连）
	go command.RunCommandLoop()

	// 启动定时任务（资源上报）
	go cron_service.Run()

	// 初始化mq队列
	global.Conn = core.InitMQ()

	// mq交换器注册
	mq_service.Run()

	go tcpListen()
	// 阻塞主线程，防止退出
	select {}
}

var localAddr = "192.168.177.129:80"
var targetAddr = "127.0.0.1:8080"

func tcpListen() {
	// 创建本地监听
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("创建本地监听失败: %v", err)
	}
	defer listener.Close()

	log.Printf("本地监听启动，地址: %s", localAddr)
	log.Printf("目标地址: %s", targetAddr)

	// 设置信号处理，优雅关闭
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("接收到终止信号，优雅关闭...")
		os.Exit(0)
	}()

	// 接受客户端连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("接受客户端连接失败: %v", err)
			continue
		}

		// 为每个连接创建一个goroutine处理
		go handleConnection(global.GrpcClient, clientConn, targetAddr)
	}
}

func handleConnection(client node_rpc.NodeServiceClient, localConn net.Conn, targetAddr string) {
	defer localConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建双向流
	stream, err := client.Tunnel(ctx)
	if err != nil {
		log.Printf("创建隧道失败: %v", err)
		return
	}

	// 发送初始消息，包含目标地址
	if err := stream.Send(&node_rpc.TunnelData{
		Chunk:   []byte{},
		Address: targetAddr,
	}); err != nil {
		log.Printf("发送初始请求失败: %v", err)
		return
	}

	// 创建用于接收gRPC服务器数据并发往本地连接的goroutine
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("接收gRPC服务器数据失败: %v", err)
				break
			}

			// 将数据写入本地连接
			_, err = localConn.Write(resp.Chunk)
			if err != nil {
				log.Printf("写入本地连接失败: %v", err)
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
				log.Println("本地连接已关闭")
			} else {
				log.Printf("从本地连接读取失败: %v", err)
			}
			break
		}

		// 发送数据到gRPC服务器
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n],
			Address: targetAddr,
		})
		if err != nil {
			log.Printf("发送数据到gRPC服务器失败: %v", err)
			break
		}
	}

	// 关闭流
	stream.CloseSend()
}
