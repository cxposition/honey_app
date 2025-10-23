package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"os"
	"time"
)

func main() {
	global.Config = core.ReadConfig()
	global.Log = core.GetLogger()
	logrus.Infof("启动成功，版本：%s, 提交：%s", global.Version, global.Commit)
	global.GrpcClient = core.GetGrpcClient()

	err := register()
	if err != nil {
		logrus.Errorf("节点注册失败：%v", err)
		return
	}

	cron_service.Run()
	select {}

}

func register() error {
	if global.Config.System.Uid == "" {
		global.Config.System.Uid = uuid.NewString()
		core.SetConfig(global.Config)
	}

	// 初始化客户端
	_ip, mac, err := ip.GetNetworkInfo(global.Config.System.Network)
	if err != nil {
		logrus.Errorf("获取网络信息失败：%v", err)
		return err
	}

	// 拿到主机名
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Errorf("获取主机名失败：%v", err)
		return err
	}

	// 获取系统信息
	systemInfo, err := info.GetSystemInfo()
	if err != nil {
		logrus.Errorf("获取系统信息失败：%v", err)
		return err
	}

	logrus.Infof("系统信息：%+v", systemInfo)

	var networkList []*node_rpc.NetworkInfoMessage
	_networkList, err := info.GetNetworkList("hy-")
	if err != nil {
		logrus.Errorf("获取网络信息失败：%v", err)
		return err
	}

	for _, networkInfo := range _networkList {
		networkList = append(networkList, &node_rpc.NetworkInfoMessage{
			Network: networkInfo.Network,
			Ip:      networkInfo.Ip,
			Net:     networkInfo.Net,
			Mask:    int32(networkInfo.Mask),
		})
	}

	// 节点注册
	_, err = global.GrpcClient.Register(context.Background(), &node_rpc.RegisterRequest{
		Ip:      _ip,
		Mac:     mac,
		NodeUid: global.Config.System.Uid,
		Version: global.Version,
		Commit:  global.Commit,
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName:            hostName,
			DistributionVersion: systemInfo.Distribution,
			CoreVersion:         systemInfo.KernelVersion,
			SystemType:          systemInfo.Arch,
			StartTime:           systemInfo.BootTime.Format(time.DateTime),
		},
		NetworkList: networkList,
	})
	if err != nil {
		logrus.Errorf("节点注册失败：%v", err)
		return err
	}
	logrus.Infof("节点注册成功")
	return nil
}
