package node_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

type UpdateRequest struct {
	ID    uint   `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}

func (NodeApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	err = global.DB.Model(&model).Update("title", cr.Title).Error
	if err != nil {
		res.FailWithMsg("更新失败", c)
		return
	}

	res.OkWithMsg("节点修改成功", c)
}
