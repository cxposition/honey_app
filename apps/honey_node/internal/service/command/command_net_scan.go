package command

import (
	"fmt"
	"github.com/j-keck/arping"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/ip"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func HandleNetScan(request *node_rpc.CmdRequest, respChan chan *node_rpc.CmdResponse) {
	req := request.GetNetScanInMessage()
	fmt.Printf("网络扫描 %v\n", req)

	ipList, err := ip.ParseIPRange(req.IpRange)
	if err != nil {
		respChan <- &node_rpc.CmdResponse{
			CmdType: node_rpc.CmdType_cmdNetScanType,
			TaskID:  request.TaskID,
			NodeID:  global.Config.System.Uid,
			NetScanOutMessage: &node_rpc.NetScanOutMessage{
				End:      true,
				Progress: 0,
				Ip:       "",
			},
		}
	}

	t1 := time.Now()
	fmt.Println("开始时间:", time.Now().Format(time.DateTime))

	filterIPList := map[string]struct{}{}
	for _, s := range req.FilterIPList {
		filterIPList[s] = struct{}{}
	}

	iface := req.Network
	concurrency := 10
	taskNum := make(chan struct{}, concurrency) // 并发限制

	var wg sync.WaitGroup
	var doneCount int64
	total := int64(len(ipList))

	for _, _ip := range ipList {
		if _, exists := filterIPList[_ip]; exists {
			continue
		}

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

			// 根据mac地址查询设备产商名称
			manufName := core.GetVendorByMAC(mac.String())

			atomic.AddInt64(&doneCount, 1)
			done := atomic.LoadInt64(&doneCount)
			progress := float64(done) / float64(total) * 100

			respChan <- &node_rpc.CmdResponse{
				CmdType: node_rpc.CmdType_cmdNetScanType,
				TaskID:  request.TaskID,
				NodeID:  global.Config.System.Uid,
				NetScanOutMessage: &node_rpc.NetScanOutMessage{
					End:      false,
					Progress: float32(progress),
					NetID:    req.NetID,
					Ip:       _ip,
					Mac:      mac.String(),
					Manuf:    manufName,
				},
			}
		}(_ip)
	}

	wg.Wait()

	respChan <- &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  request.TaskID,
		NodeID:  global.Config.System.Uid,
		NetScanOutMessage: &node_rpc.NetScanOutMessage{
			End: true,
		},
	}

	fmt.Printf("\n扫描完成，总耗时: %v\n", time.Since(t1))
	fmt.Println("结束时间:", time.Now().Format(time.DateTime))
}
