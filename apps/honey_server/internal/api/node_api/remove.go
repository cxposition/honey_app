package node_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

func (NodeApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 判断这个节点是否在运行
	// 可以通过status判断，但是它的实时性不太准
	err = global.DB.Delete(&model).Error
	if err != nil {
		res.FailWithMsg("节点删除失败", c)
		return
	}

	res.OkWithMsg("节点删除成功", c)
}
