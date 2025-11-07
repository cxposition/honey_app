package honey_ip_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
)

func (HoneyIPApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequestList](c)
	var honeyIPList []models.HoneyIpModel
	global.DB.Preload("NodeModel").Find(&honeyIPList, "id in ?", cr.IdList)
	if len(honeyIPList) == 0 {
		res.FailWithMsg("未找到诱捕ip", c)
		return
	}

	nodeModel := honeyIPList[0].NodeModel
	// 判断节点是否在线
	if nodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 使用封装的获取节点函数
	_, ok := grpc_service.GetNodeCommand(nodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 下发批量删除的任务到消息队列

	// 改状态
	global.DB.Model(&honeyIPList).Update("status", 4)
	res.OkWithMsg("批量删除中", c)
}
