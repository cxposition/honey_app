package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Options struct {
	NodeID string
}

var options Options

func main() {
	flag.StringVar(&options.NodeID, "c", "", "节点id")
	flag.Parse()

	fmt.Println(options.NodeID)

	// 连接RabbitMQ
	conn, err := amqp.Dial("amqp://admin:password@192.168.177.129:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明与生产者一致的交换器（确保交换器存在）
	err = ch.ExchangeDeclare(
		"create_ip_queue", // 交换器名称（与生产者一致）
		"direct",          // 交换器类型（直接交换器）
		true,              // 持久化
		false,             // 自动删除
		false,             // 内部交换器
		false,             // 非阻塞
		nil,               // 参数
	)
	failOnError(err, "Failed to declare create_ip_queue exchange")

	// --------------------------
	// 为消费者1创建队列并绑定（node01）
	// --------------------------
	queue1, err := ch.QueueDeclare(
		fmt.Sprintf("create_ip_%s_queue", options.NodeID), // 队列名称（唯一标识，与node01绑定）
		true,  // 持久化队列
		false, // 不自动删除
		false, // 非排他性
		false, // 非阻塞
		nil,   // 参数
	)
	failOnError(err, "Failed to declare node01 queue")

	// 绑定队列到交换器，绑定键为 "node01"（只接收路由键为node01的消息）
	err = ch.QueueBind(
		queue1.Name,       // 队列名称
		options.NodeID,    // 绑定键（与生产者路由键匹配）
		"create_ip_queue", // 交换器名称
		false,             // 非阻塞
		nil,               // 参数
	)
	failOnError(err, "Failed to bind node01 queue")

	// --------------------------
	consumeMessages(ch, queue1.Name) // 消费者1处理node01消息
}

// 通用消息消费函数（接收队列名称和消费者标识）
func consumeMessages(ch *amqp.Channel, queueName string) {
	// 从队列接收消息
	msgs, err := ch.Consume(
		queueName, // 队列名称（从哪个队列消费）
		"",        // 消费者标识（自定义名称）
		false,     // 关闭自动确认（手动确认消息处理完成）
		false,     // 非排他性
		false,     // 非本地消费者
		false,     // 非阻塞
		nil,       // 参数
	)
	failOnError(err, "Failed to register consumer")

	// 循环处理消息
	for d := range msgs {
		var ip string
		// 解析消息内容（生产者发送的是IP字符串）
		if err := json.Unmarshal(d.Body, &ip); err != nil {
			fmt.Println(err)
			d.Nack(false, false) // 拒绝消息，不重新入队
			continue
		}

		// 处理消息（模拟创建IP的业务逻辑）
		fmt.Printf("Received task: Create IP %s (Routing Key: %s)\n", ip, d.RoutingKey)
		// TODO: 实际业务逻辑（如调用创建IP的API）

		// 手动确认消息已处理完成（消息从队列中删除）
		if err := d.Ack(false); err != nil {
			fmt.Printf("Error acknowledging message: %v\n", err)
		} else {
			fmt.Printf("Completed task: IP %s created\n", ip)
		}
	}
}
