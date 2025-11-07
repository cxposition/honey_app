package node_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

func (NodeApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[models.PageInfo](c)
	list, count, _ := common_service.QueryList(models.NodeModel{}, common_service.ListRequest{
		Likes:    []string{"title", "ip"}, // username like req.Key
		PageInfo: cr,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
