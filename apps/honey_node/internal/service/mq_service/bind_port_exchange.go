package mq_service

import "fmt"

func BindPortExchange(msg string) error {
	fmt.Println("绑定端口消息", msg)
	return nil
}
