package command

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
)

func HandleNetworkFlush(req *node_rpc.CmdRequest, respChan chan *node_rpc.CmdResponse) {
	if len(req.NetworkFlushInMessage.FilterNetworkName) == 0 {
		logrus.Warn("⚠️ 命令参数为空")
		return
	}
	name := req.NetworkFlushInMessage.FilterNetworkName[0]

	list, err := info.GetNetworkList([]string{name})
	if err != nil {
		logrus.Errorf("获取网络信息失败：%v", err)
		return
	}

	var outList []*node_rpc.NetworkInfoMessage
	for _, n := range list {
		outList = append(outList, &node_rpc.NetworkInfoMessage{
			Network: n.Network,
			Ip:      n.Ip,
			Net:     n.Net,
			Mask:    int32(n.Mask),
		})
	}

	resp := &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  req.TaskID,
		NodeID:  global.Config.System.Uid,
		NetworkFlushOutMessage: &node_rpc.NetworkFlushOutMessage{
			NetworkList: outList,
		},
	}

	select {
	case respChan <- resp:
	default:
		logrus.Warn("⚠️ 响应通道已满，丢弃响应")
	}
}
