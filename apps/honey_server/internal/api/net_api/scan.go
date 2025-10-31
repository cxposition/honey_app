package net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
	"time"
)

func (NetApi) ScanView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NetModel
	if err := global.DB.Preload("NodeModel").Take(&model, cr.ID).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	if model.NodeModel.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	// 3️⃣ 构造命令请求（让节点刷新网卡）
	taskID := uuid.New().String()
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  taskID,
		NetScanInMessage: &node_rpc.NetScanInMessage{
			Network:      model.Network,
			IpRange:      model.CanUseHoneyIPRange,
			FilterIPList: []string{},
			NetID:        uint32(model.ID),
		},
	}

	// 4️⃣ 通过 gRPC 下发命令
	resp, err := grpc_service.SendCommand(model.NodeModel.Uid, req, 8*time.Second)
	if err != nil {
		res.FailWithMsg(fmt.Sprintf("命令执行失败: %v", err), c)
		return
	}
	if resp == nil || resp.NetScanOutMessage == nil {
		res.FailWithMsg("节点未返回网卡信息", c)
		return
	}

	logrus.Infof("收到节点的网卡扫描返回消息%+v", resp)
}
