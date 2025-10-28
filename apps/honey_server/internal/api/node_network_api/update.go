package node_network_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

type UpdateRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Gateway string `json:"gateway"`
}

func (NodeNetworkApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)
	var model models.NodeNetworkModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点网卡不存在", c)
		return
	}

	// 需要判断这个网关
	// 1.必须得是ipv4
	// 2.这个ip不能是探针ip

	err = global.DB.Model(&model).Update("gateway", cr.Gateway).Error
	if err != nil {
		res.FailWithMsg("网卡修改失败", c)
		return
	}
	res.FailWithMsg("节点网卡修改成功", c)
}
