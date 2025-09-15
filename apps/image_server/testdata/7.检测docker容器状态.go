package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerClient 是全局客户端实例（建议单例）
var dockerClient *client.Client

// InitDockerClient 初始化 Docker 客户端
func InitDockerClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("无法初始化 Docker 客户端: %w", err)
	}
	dockerClient = cli
	return nil
}

// ListAllContainers 获取所有容器的状态（包括运行中、已停止等）
func ListAllContainers() ([]container.Summary, error) {
	// 创建Docker客户端
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("创建Docker客户端失败: %v", err)
	}
	defer cli.Close()

	// 获取所有容器列表（包括停止的容器）
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	return containers, nil
}

// GetContainerStatusByName 根据容器名获取其状态
func GetContainerStatusByName(containerName string) (container.Summary, error) {
	// 创建Docker客户端
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return container.Summary{}, fmt.Errorf("创建Docker客户端失败: %v", err)
	}
	defer cli.Close()

	// 使用过滤器按名称查找容器
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: filter,
		All:     true,
	})
	if err != nil {
		return container.Summary{}, fmt.Errorf("获取容器列表失败: %v", err)
	}

	if len(containers) == 0 {
		return container.Summary{}, fmt.Errorf("未找到名为 %s 的容器", containerName)
	}

	// 返回第一个匹配的容器（容器名应该是唯一的）
	return containers[0], nil
}

// 示例主函数（用于测试）
func main() {
	// 初始化 Docker 客户端
	if err := InitDockerClient(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 所有容器状态 ===")
	allContainers, err := ListAllContainers()
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range allContainers {
		fmt.Printf("ID: %s, 名称: %s, 状态: %s\n", container.ID[:12], container.Names[0][1:], container.State)
	}

	fmt.Println("\n=== 单个容器状态 ===")
	containerName := "redis" // 替换为你实际的容器名
	singleContainer, err := GetContainerStatusByName(containerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n容器 %s 的状态: %s\n", containerName, singleContainer.State)
}
