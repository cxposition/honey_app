package mq_service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_server/internal/global"
)

type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIPID"`
	IP        string `json:"ip"`
	Mask      int8   `json:"mask"`
	Network   string `json:"network"` // 基于哪个接口创建
	IsTan     bool   `json:"isTan"`   // 是否是探针ip
	LogID     string `json:"logID"`   // 日志id
}

func SendCeateIPMsg(nodeUID string, req CreateIPRequest) {
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ
	// 发送消息

	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Errorf("创建channel失败 %s", err)
		return
	}
	err = ch.Publish(cfg.CreateIpExchangeName,
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
