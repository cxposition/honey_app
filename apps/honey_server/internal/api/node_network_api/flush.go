package node_network_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
)

func (NodeNetworkApi) FlushView(c *gin.Context) {
	grpc_service.CmdRequestChan <- &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  "xxx",
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}
	res.OkWithMsg("刷新成功", c)
}
