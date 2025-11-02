package main

import (
	"fmt"
	"github.com/j-keck/arping"
	"honey_node/internal/utils/ip"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	t1 := time.Now()
	fmt.Println("开始时间:", time.Now().Format(time.DateTime))

	ipList, err := ip.ParseIPRange("192.168.177.1-192.168.177.200")
	if err != nil {
		fmt.Println("解析IP段错误:", err)
		return
	}

	iface := "ens33"
	taskNum := make(chan struct{}, 100) // 并发限制

	var wg sync.WaitGroup
	var doneCount int64
	total := int64(len(ipList))

	for _, _ip := range ipList {
		taskNum <- struct{}{}
		wg.Add(1)

		go func(ipStr string) {
			defer wg.Done()
			defer func() { <-taskNum }()

			mac, duration, err := arping.PingOverIfaceByName(net.ParseIP(ipStr), iface)
			if err != nil {
				fmt.Printf("\n[发现设备] IP: %-15s MAC: %s RTT: %v\n", ipStr, mac, duration)
			} else {
				fmt.Printf("\n[发现设备] IP: %-15s MAC: %s RTT: %v\n", ipStr, mac, duration)
			}

			atomic.AddInt64(&doneCount, 1)
			done := atomic.LoadInt64(&doneCount)
			fmt.Printf("\r进度: %.1f%% (%d/%d)", float64(done)/float64(total)*100, done, total)
		}(_ip)
	}

	wg.Wait()

	fmt.Printf("\n扫描完成，总耗时: %v\n", time.Since(t1))
	fmt.Println("结束时间:", time.Now().Format(time.DateTime))
}
