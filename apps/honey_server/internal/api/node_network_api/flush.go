package node_network_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
)

func (NodeNetworkApi) FlushView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	if model.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	_, ok := grpc_service.NodeCommandMap[model.Uid]
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	grpc_service.NodeCommandMap[model.Uid].ReqChan <- &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  "xxx",
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	// 拿到节点的数据
	response := <-grpc_service.NodeCommandMap[model.Uid].ResChan
	fmt.Println("网卡刷新数据", response)
	res.OkWithData(response.NetworkFlushOutMessage, c)
}
