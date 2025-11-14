package mq_service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_node/internal/global"
)

type AlertMsgType struct {
	NodeUid          string `json:"nodeUid"`
	SrcIp            string `json:"srcIp"`
	SrcPort          int    `json:"srcPort"`
	DestIP           string `json:"destIP"`
	DestPort         int    `json:"destPort"`
	Timestamp        string `json:"timestamp"` // 年月日，时分秒的时间
	Signature        string `json:"signature"`
	Level            int8   `json:"level"` // 告警级别
	HttpResponseBody string `json:"httpResponseBody"`
	Payload          string `json:"payload"`
}

func SendAlertMsg(data AlertMsgType) error {
	cfg := global.Config.MQ

	// 声明队列
	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
		return err
	}
	_, err = ch.QueueDeclare(
		cfg.AlertTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalf("Failed to declare a queue: %v", err)
		return err
	}
	// 发送消息
	byteData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("SendAlert json.Marshal err: %v", err)
		return err
	}
	err = ch.Publish("",
		cfg.AlertTopic,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byteData,
		})
	if err != nil {
		logrus.Errorf("发送告警信息失败: %v %s", err, data)
		return err
	}
	logrus.Infof("发送告警信息成功: %+v", data)
	return nil
}
