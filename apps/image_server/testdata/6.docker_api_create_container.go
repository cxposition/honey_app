package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

func main() {
	// ================== 可配置参数（你只需修改下面这些）==================

	const (
		// 镜像名称（支持 tag）
		ImageName = "nginx:latest"

		// 容器名称（必须唯一）
		ContainerName = "my-nginx"

		// 自定义网络名称（必须是用户自定义网络，不能是 default bridge）
		NetworkName = "honey-hy"

		// 指定容器在该网络中的静态 IP（必须属于网络子网范围）
		ContainerIP = "10.2.0.15"

		// 主机端口映射：主机端口 -> 容器端口
		HostPort      = "8888"
		ContainerPort = "80"
	)

	// ====================================================================

	// 初始化 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("无法初始化 Docker 客户端:", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// 配置容器镜像和命令
	containerConfig := &container.Config{
		Image: ImageName,
		Cmd:   []string{"nginx", "-g", "daemon off;"},
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(ContainerPort + "/tcp"): {},
		},
	}

	// 端口绑定配置
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			nat.Port(ContainerPort + "/tcp"): {{
				HostIP:   "0.0.0.0",
				HostPort: HostPort,
			}},
		},
	}

	// 网络配置：指定网络和静态 IP
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			NetworkName: {
				IPAMConfig: &network.EndpointIPAMConfig{ // ← 注意：是 EndpointIPAMConfig
					IPv4Address: ContainerIP, // ✅ 这里才是正确位置
				},
			},
		},
	}

	// 创建容器（使用指定名称）
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, ContainerName)
	if err != nil {
		log.Fatal("创建容器失败:", err)
	}

	fmt.Printf("✅ 容器创建成功，ID: %s\n", resp.ID)

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatal("启动容器失败:", err)
	}

	fmt.Printf("🚀 容器 %s 已启动，访问 http://localhost:%s 查看服务\n", ContainerName, HostPort)
}
