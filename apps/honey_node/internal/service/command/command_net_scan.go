package command

import (
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
)

func HandleNetScan(request *node_rpc.CmdRequest, respChan chan *node_rpc.CmdResponse) {
	request.GetNetScanInMessage()
	fmt.Printf("网络扫描 %v\n", request)

	respChan <- &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  request.TaskID,
		NodeID:  global.Config.System.Uid,
		NetScanOutMessage: &node_rpc.NetScanOutMessage{
			End:      false,
			Progress: 0,
			Ip:       "192.168.100.1",
		},
	}

	//time.Sleep(2 * time.Second)
	//
	//respChan <- &node_rpc.CmdResponse{
	//	CmdType: node_rpc.CmdType_cmdNetScanType,
	//	TaskID:  request.TaskID,
	//	NodeID:  global.Config.System.Uid,
	//	NetScanOutMessage: &node_rpc.NetScanOutMessage{
	//		End:      true,
	//		Progress: 0,
	//		Ip:       "192.168.100.1",
	//	},
	//}

}
