package node_network_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

type ListRequest struct {
	NodeID uint `form:"nodeID" binding:"required"`
	models.PageInfo
}

func (NodeNetworkApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	list, count, _ := common_service.QueryList(models.NodeNetworkModel{NodeID: cr.NodeID}, common_service.ListRequest{
		Likes:    []string{"network", "ip"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
