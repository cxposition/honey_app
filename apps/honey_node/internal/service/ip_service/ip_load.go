package ip_service

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/models"
	"honey_node/internal/utils"
	"honey_node/internal/utils/info"
)

func IPLoad() {
	var ipList []models.IpModel
	global.DB.Find(&ipList)
	networkMap, err := info.GetNetworkInterfaces()
	if err != nil {
		logrus.Errorf("获取网卡信息失败: %s", err)
		return
	}

	for _, model := range ipList {
		// 判断ip在不在这个机器上，如果在就不创建了
		ips, ok := networkMap[model.LinkName]
		if ok {
			// 在里面，判断ip地址对不对
			if !utils.Inlist(ips, model.Ip) {
				logrus.Errorf("网卡%s对应的ip地址错误 %v %s", model.LinkName, ips, model.Ip)
			}
			continue
		}
		// 创建网卡和ip
		_, err = SetIp(SetIpRequest{
			Ip:       model.Ip,
			Mask:     model.Mask,
			LinkName: model.LinkName,
			Network:  model.Network,
			Mac:      model.Mac,
		})
		if err != nil {
			logrus.Errorf("初始化ip错误: %s", err)
			continue
		}
	}
}
