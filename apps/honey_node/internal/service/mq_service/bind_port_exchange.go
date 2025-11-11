package mq_service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_node/internal/service/port_service"
)

type BindPortRequest struct {
	IP       string     `json:"ip"`
	PortList []PortInfo `json:"portList"`
	LogID    string     `json:"logID"` // 日志id
}

type PortInfo struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DestIP   string `json:"destIP"`
	DestPort int    `json:"destPort"`
}

func (p PortInfo) LocalAddr() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (p PortInfo) TargetAddr() string {
	return fmt.Sprintf("%s:%d", p.DestIP, p.DestPort)
}

func BindPortExchange(msg string) error {
	logrus.Infof("绑定端口消息:%#v", msg)

	var req BindPortRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("JSON 解析失败:%v, 消息%s", err, msg)
		return nil
	}

	// 先把之前这个ip上的服务全部停止
	port_service.CloseIpTunnel(req.IP)

	for _, port := range req.PortList {
		// 起端口监听
		go func(port PortInfo) {
			err := port_service.Tunnel(port.LocalAddr(), port.TargetAddr())
			if err != nil {
				logrus.Errorf("端口绑定失败 %s", err)
			}
			// 如果报错，大概率是ip没有起来，也有可能是端口没有释放掉
			// 需要通知管理，只通知失败的，成功的不通知

		}(port)
	}

	return nil
}
