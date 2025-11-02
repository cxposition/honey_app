package net_api

import (
	"context"
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

	nodeID := model.NodeModel.Uid
	val, ok := grpc_service.NodeCommandMap.Load(nodeID)
	if !ok {
		logrus.Errorf("节点 %s 未找到", nodeID)
		//cancelFunc()
		return
	}
	cmd := val.(*grpc_service.Command)

	respChan := make(chan *node_rpc.CmdResponse, 1)
	cmd.ResMap.Store(req.TaskID, respChan)
	defer cmd.ResMap.Delete(req.TaskID)

	select {
	case cmd.ReqChan <- req:
		// 成功发送
		logrus.Infof("节点 %s 的 task %s 已发送", nodeID, req.TaskID)
	case <-time.After(3 * time.Second):
		logrus.Errorf("发送命令到节点 %s 超时", nodeID)
	}

	// ✅ 添加超时机制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

label:
	for {
		select {
		case resp := <-respChan:
			logrus.Debugf("节点 %s 的 task %s 的结果为 %+v", nodeID, req.TaskID, resp)
			message := resp.NetScanOutMessage
			logrus.Infof("节点信息为:%+v", message)
			if message.ErrMsg != "" {
				res.FailWithMsg("扫描错误"+message.ErrMsg, c)
				break label
			}
			if message.End {
				break label
			}
		case <-ctx.Done():
			res.FailWithMsg("获取响应超时", c)
			return
		default:
		}
	}

	res.OkWithMsg("扫描成功", c)

}
