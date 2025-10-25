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
	"time"

	"github.com/google/uuid"
)

func (NodeNetworkApi) FlushView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NodeModel
	if err := global.DB.Take(&model, cr.ID).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	if model.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	taskID := uuid.New().String()
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  taskID,
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	resp, err := grpc_service.SendCommand(model.Uid, req, 5*time.Second)
	if err != nil {
		res.FailWithMsg(fmt.Sprintf("命令执行失败: %v", err), c)
		return
	}

	res.OkWithData(resp.NetworkFlushOutMessage, c)
}
