package mq_service

import (
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
)

func RegisterExchange() {
	cfg := global.Config.MQ
	exchangeDeclare(cfg.CreateIpExchangeName)
	exchangeDeclare(cfg.BindPortExchangeName)
	exchangeDeclare(cfg.DeleteIpExchangeName)
}

func exchangeDeclare(name string) {
	// 声明交换机
	err := global.Queue.ExchangeDeclare(
		name,     // 交换机名称
		"direct", // 直接交换机类型
		true,     // 持久化
		false,    //自动删除
		false,    //内部
		false,    //非阻塞
		nil,      //参数
	)
	if err != nil {
		logrus.Fatalf("Failed to declare exchange: %v", err)
	}
	logrus.Infof("声明交换器 %s 成功", name)
}
