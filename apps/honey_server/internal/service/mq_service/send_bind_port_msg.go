package mq_service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_server/internal/global"
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

func SendBindPortMsg(nodeUID string, req BindPortRequest) {
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ
	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Errorf("Failed to open a channel: %v", err)
		return
	}
	// 发送消息
	err = ch.Publish(cfg.BindPortExchangeName,
		nodeUID,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byteData,
		},
	)
	if err != nil {
		logrus.Errorf("消息发送失败 %s %s", err, string(byteData))
	} else {
		logrus.Infof("消息发送成功 %s", string(byteData))
	}
}
