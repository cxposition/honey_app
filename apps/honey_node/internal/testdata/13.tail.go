package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"io"
)

func main() {
	t, err := tail.TailFile("../../deploy/suricata/logs/eve.json",
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: io.SeekEnd, // 从最新的地方开始监听
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
		// json解析
		// 发送到mq
	}
}
