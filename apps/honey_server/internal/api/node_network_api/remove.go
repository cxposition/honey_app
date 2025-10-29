package node_network_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

func (NodeNetworkApi) RemoveView(c *gin.Context) {
	// 如果网卡没启用，则直接删除
	// 如果网卡启用了，但是没有诱捕ip，也直接删除
	// 但是网卡启用了，并且也创建了诱捕ip，则不能删除
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeNetworkModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网卡不存在", c)
		return
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		res.FailWithMsg("删除网卡失败", c)
		return
	}
	res.OkWithMsg("删除网卡成功", c)
	return
}
