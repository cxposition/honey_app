package node_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

func (NodeApi) DetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	res.OkWithData(model, c)
}
