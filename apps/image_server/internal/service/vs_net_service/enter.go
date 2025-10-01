package vs_net_service

import (
	"context"
	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
)

func Run() {
	cfg := global.Config.VsNet

	// 设置默认值
	if cfg.Name == "" {
		cfg.Name = "honey-hy"
	}
	if cfg.Net == "" {
		cfg.Net = "10.2.0.0/24"
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "hy-"
	}

	// 判断配置文件中的网络名称，在不在列表中
	// 如果不在，那就创建对应的docker网络
	// 如果在，那就判断这个docker网络对应的子网，是不是和配置文件里面的子网一样
	// 如果不一样就提示用户，需要排查问题
	// 获取所有网络
	networks, err := global.DockerClient.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		logrus.Fatalf("获取虚拟网络列表失败: %v", err)
	}

	// 查找配置中的网络是否存在
	var found bool
	var existingNetwork network.Summary
	for _, networkT := range networks {
		if networkT.Name == cfg.Name {
			found = true
			existingNetwork = networkT
			break
		}
	}

	// 处理网络不存在的情况
	if !found {
		// 创建网络
		ipam := network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet: cfg.Net,
				},
			},
		}
		_, err := global.DockerClient.NetworkCreate(context.Background(), cfg.Name, network.CreateOptions{
			IPAM: &ipam,
		})
		if err != nil {
			logrus.Fatalf("创建网络失败: %v", err)
		}
		logrus.Printf("成功创建网络 %s，子网为 %s", cfg.Name, cfg.Net)
		return
	}

	// 网络存在，检查子网是否匹配
	// 注意：这里简化了子网检查逻辑，实际实现可能需要更复杂的CIDR比较
	// 这里假设网络配置中只有一个IPAM配置
	if len(existingNetwork.IPAM.Config) > 0 && existingNetwork.IPAM.Config[0].Subnet != cfg.Net {
		logrus.Warnf("警告: 网络 %s 存在，但子网不匹配。现有子网: %s，配置子网: %s",
			cfg.Name, existingNetwork.IPAM.Config[0].Subnet, cfg.Net)
		logrus.Fatalf("请排查网络配置问题")
		return
	}

	logrus.Infof("网络 %s 存在且子网匹配: %s", cfg.Name, cfg.Net)
}
