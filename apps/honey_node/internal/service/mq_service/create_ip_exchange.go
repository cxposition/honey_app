package mq_service

import "fmt"

func CreateIpExchange(msg string) error {
	fmt.Println("消息", msg)
	return nil
}
