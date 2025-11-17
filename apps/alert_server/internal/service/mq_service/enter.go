package mq_service

import (
	"alert_server/internal/global"
	"github.com/sirupsen/logrus"
)

func Run() {
	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Fatalf("创建MQ通道失败 %s", err)
	}
	cfg := global.Config.Alert
	// 声明队列
	_, err = ch.QueueDeclare(cfg.AlertTopic, false, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("声明mq队列失败 %s", err)
	}

	go RevAlertMQ(ch)

}
