package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// 1. 加载客户端证书和密钥（双向认证时需要）
	cert, err := tls.LoadX509KeyPair("G:\\honey_app\\deploy\\rabbitMQ\\ssl\\client_certificate.pem", "G:\\honey_app\\deploy\\rabbitMQ\\ssl\\client_key.pem")
	if err != nil {
		log.Fatalf("加载客户端证书失败: %v", err)
	}

	// 2. 加载CA证书（验证服务器证书）
	caCert, err := ioutil.ReadFile("G:\\honey_app\\deploy\\rabbitMQ\\ssl\\ca_certificate.pem")
	if err != nil {
		log.Fatalf("读取CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. 配置TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert}, // 客户端证书（双向认证时需要）
		RootCAs:            caCertPool,              // 信任的CA
		InsecureSkipVerify: false,                   // 必须验证服务器证书
	}

	// 4. 创建AMQP连接（使用TLS）
	conn, err := amqp.DialTLS("amqps://admin:password@192.168.177.129:5671/", tlsConfig)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 其余代码与普通连接相同...

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

	// 发送消息
	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}
	fmt.Printf("发送消息: %s\n", body)
}
