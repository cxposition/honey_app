package docker_service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

func RunContainer(containerName, networkName, ip, imageName string) (containerID string, err error) {
	// 初始化 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("无法初始化 Docker 客户端:", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// 配置容器镜像和命令
	containerConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: map[nat.Port]struct{}{
			//nat.Port(ContainerPort + "/tcp"): {},
		},
	}

	// 端口绑定配置
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			//nat.Port(ContainerPort + "/tcp"): {{
			//	HostIP: "0.0.0.0",
			//	HostPort: HostPort,
			//}},
		},
	}

	// 网络配置：指定网络和静态 IP
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				IPAMConfig: &network.EndpointIPAMConfig{ // ← 注意：是 EndpointIPAMConfig
					IPv4Address: ip, // ✅ 这里才是正确位置
				},
			},
		},
	}

	// 创建容器（使用指定名称）
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		log.Fatal("创建容器失败:", err)
	}

	fmt.Printf("✅ 容器创建成功，ID: %s\n", resp.ID)

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Fatal("启动容器失败:", err)
	}

	//fmt.Printf("🚀 容器 %s 已启动，访问 http://localhost:%s 查看服务\n", containerName, HostPort)
	return
}
