package docker_service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"image_server/internal/global"
)

func ListAllContainers() ([]container.Summary, error) {
	// 创建Docker客户端
	ctx := context.Background()
	// 获取所有容器列表（包括停止的容器）
	containers, err := global.DockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	return containers, nil
}

// GetContainerStatus 根据容器名获取其状态
func GetContainerStatus(containerName string) (container.Summary, error) {
	ctx := context.Background()
	// 使用过滤器按名称查找容器
	filter := filters.NewArgs()
	filter.Add("name", containerName)
	containers, err := global.DockerClient.ContainerList(ctx, container.ListOptions{
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

// PrefixContainerStatus 根据容器名前缀获取容器状态
func PrefixContainerStatus(containerName string) (summaryList []container.Summary, err error) {
	// 使用过滤器按名称查找容器
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	containers, err := global.DockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: filter,
		All:     true,
	})
	if err != nil {
		return
	}
	// 返回第一个匹配的容器（容器名应该是唯一的）
	return containers, nil
}
