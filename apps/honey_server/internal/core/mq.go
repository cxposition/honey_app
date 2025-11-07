package core

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_server/internal/global"
	"log"
	"os"
)

func InitMQ() *amqp.Channel {
	cfg := global.Config.MQ
	var conn *amqp.Connection
	var err error
	if cfg.Ssl {
		// 1. 加载客户端证书和密钥（双向认证时需要）
		cert, err := tls.LoadX509KeyPair(cfg.ClientCertificate, cfg.ClientKey)
		if err != nil {
			log.Fatalf("加载客户端证书失败: %v", err)
		}

		// 2. 加载CA证书（验证服务器证书）
		caCert, err := os.ReadFile(cfg.CaCertificate)
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
		conn, err = amqp.DialTLS(cfg.Addr(), tlsConfig)
	} else {
		conn, err = amqp.Dial(cfg.Addr())
	}

	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a channel: %v", err)
	}
	return ch
}

func RegisterExchange() {}
