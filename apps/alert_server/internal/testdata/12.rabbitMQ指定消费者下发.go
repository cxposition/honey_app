package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func init() {
	_conn, err := amqp.Dial("amqp://admin:password@192.168.177.129:5672/")
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	conn = _conn
}

var ch *amqp.Channel

func init() {
	// 创建通道
	_ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	ch = _ch
}
func InitExchangeDeclare() (err error) {
	// 声明交换器
	err = ch.ExchangeDeclare(
		"create_ip_queue", // 交换器名称
		"direct",          // 直接交换器类型
		true,              // 持久化
		false,             // 自动删除
		false,             // 内部
		false,             // 非阻塞
		nil,               // 参数
	)
	if err != nil {
		return err
	}

	return err

}

func SendMsg(key string, nodeUID string, val any) (err error) {
	byteData, _ := json.Marshal(val)
	// 发送消息
	err = ch.Publish(
		key,     // exchange
		nodeUID, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byteData,
		})
	return err
}

func main() {
	err := InitExchangeDeclare()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(SendMsg("create_ip_queue", "node01", "192.168.100.1"))
	fmt.Println(SendMsg("create_ip_queue", "node02", "192.168.200.1"))
}
