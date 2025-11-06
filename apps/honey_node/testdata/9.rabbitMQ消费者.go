package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// 连接 RabbitMQ
	conn, err := amqp.Dial("amqp://admin:password@192.168.177.129:5672/")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		"hello", // 队列名称
		false,   // 持久性
		false,   // 自动删除
		false,   // 排他性
		false,   // 非阻塞
		nil,     // 其他参数
	)
	if err != nil {
		log.Fatalf("无法声明队列: %v", err)
	}

	// 接收消息
	msgs, err := ch.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 其他参数
	)
	if err != nil {
		log.Fatalf("无法注册消费者: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("收到消息: %s\n", d.Body)
		}
	}()

	fmt.Println("等待消息中... 按 CTRL+C 退出")
	<-forever
}
