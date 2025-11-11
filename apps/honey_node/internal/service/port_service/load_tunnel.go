package port_service

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/models"
)

func LoadTunnel() {
	var portList []models.PortModel
	global.DB.Find(&portList)
	logrus.Infof("加载端口转发记录 %d", len(portList))
	for _, model := range portList {
		go Tunnel(model.LocalAddr, model.TargetAddr)
	}
}
