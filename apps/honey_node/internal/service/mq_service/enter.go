package mq_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
)

func Run() {
	cfg := global.Config.MQ
	// 由于接收消息会阻塞，这里采用异步方式启动交换器注册
	go register(cfg.CreateIpExchangeName, CreateIpExchange)
	go register(cfg.DeleteIpExchangeName, DeleteIpExchange)
	go register(cfg.BindPortExchangeName, BindPortExchange)
	logrus.Infof("启动MQ服务成功")
}

func register(exchangeName string, fn func(msg string) error) {
	// 声明与生产者一致的交换器（确保交换器存在）
	ch, err := global.MQConn.Channel()
	if err != nil {
		logrus.Fatalf("[%s] 创建通道失败: %v", exchangeName, err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchangeName, // 交换器名称（与生产者一致）
		"direct",     // 交换器类型（直接交换器）
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部交换器
		false,        // 非阻塞
		nil,          // 参数
	)
	if err != nil {
		logrus.Fatalf("%s 声明交换器失败: %v", exchangeName, err)
	}

	_cfg := global.Config
	// 为消费者创建队列并绑定
	queue, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s_queue", exchangeName, _cfg.System.Uid), // 队列名称（唯一标识，与node01绑定）
		true,  // 持久化队列
		false, // 不自动删除
		false, // 非排他性
		false, // 非阻塞
		nil,   // 参数
	)
	if err != nil {
		logrus.Fatalf("%s 创建队列失败: %v", exchangeName, err)
	}

	// 绑定队列到交换器，绑定键为节点id
	err = ch.QueueBind(
		queue.Name,      // 队列名称
		_cfg.System.Uid, // 绑定键（与生产者路由键匹配）
		exchangeName,    // 交换器名称
		false,           // 非阻塞
		nil,             // 参数
	)
	if err != nil {
		logrus.Fatalf("%s 绑定队列失败: %v", exchangeName, err)
	}

	// 从队列接收消息
	_megs, err := ch.Consume(
		queue.Name, // 队列名称（从哪个队列消费）
		"",         // 消费者标识（自定义名称）
		false,      // 关闭自动确认（手动确认消息处理完成）
		false,      // 非排他性
		false,      // 非本地消费者
		false,      // 非阻塞
		nil,        // 参数
	)

	// 循环处理消息
	for d := range _megs {
		err = fn(string(d.Body)) // 回调函数处理消息
		if err != nil {
			d.Nack(false, true) // 拒接消息, 重新入队
			continue
		}
		d.Ack(false) // 确认收到消息，不重新发送
	}
}
