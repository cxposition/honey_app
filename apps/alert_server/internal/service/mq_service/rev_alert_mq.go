package mq_service

import (
	"alert_server/internal/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func RevAlertMQ(ch *amqp.Channel) {
	cfg := global.Config.Alert
	msgs, err := ch.Consume(cfg.AlertTopic, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("消费mq队列失败 %s", err)
	}
	if err != nil {
		logrus.Fatalf("消费mq队列失败 %s", err)
	}

	for d := range msgs {
		fmt.Printf("接收到消息 %s", d.Body)
	}
}
