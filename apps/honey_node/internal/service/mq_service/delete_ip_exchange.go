package mq_service

import "fmt"

func DeleteIpExchange(msg string) error {
	fmt.Println("收到删除消息", msg)
	return nil
}
