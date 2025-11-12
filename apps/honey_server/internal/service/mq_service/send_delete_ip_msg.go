package mq_service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_server/internal/global"
)

type DeleteIPRequest struct {
	IpList []IpInfo `json:"ipList"`
	LogID  string   `json:"logID"` // 日志id
}

type IpInfo struct {
	HoneyIPID uint   `json:"honeyIPID"`
	IP        string `json:"ip"`
	Network   string `json:"network"` // 基于哪个接口创建
	IsTan     bool   `json:"isTan"`   // 是否是探针ip
}

func SendDeleteIPMsg(nodeUID string, req DeleteIPRequest) {
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ
	// 发送消息

	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Errorf("消息发送失败 %s %s", err, string(byteData))
		return
	}
	err = ch.Publish(cfg.DeleteIpExchangeName,
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
