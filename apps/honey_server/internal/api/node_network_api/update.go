package node_network_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
	"net"
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
	// 3.这个ip必须属于这个网络
	if cr.Gateway != "" {
		gateway := net.ParseIP(cr.Gateway)
		if gateway == nil || gateway.To4() == nil {
			res.FailWithMsg("网关ip格式错误", c)
			return
		}
		to4 := gateway.To4()
		if to4 == nil {
			res.FailWithMsg("网关ip只支持ipv4", c)
			return
		}
		if cr.Gateway == model.IP {
			res.FailWithMsg("网关ip不能是探针ip", c)
			return
		}
		_, _net, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", model.IP, model.Mask))
		if !_net.Contains(gateway) {
			res.FailWithMsg("网关ip不属于当前子网", c)
			return
		}
	}

	err = global.DB.Model(&model).Update("gateway", cr.Gateway).Error
	if err != nil {
		res.FailWithMsg("网卡修改失败", c)
		return
	}
	res.FailWithMsg("节点网卡修改成功", c)
}
