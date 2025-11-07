package net_api

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

type ListResponse struct {
	models.NetModel
	NodeTitle  string `json:"nodeTitle"`
	NodeStatus int8   `json:"nodeStatus"`
}

func (NetApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	_list, count, _ := common_service.QueryList(models.NetModel{NodeID: cr.NodeID}, common_service.ListRequest{
		Likes:    []string{"title", "ip"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
		Preload:  []string{"NodeModel"},
	})
	var list = make([]ListResponse, 0, len(_list))
	for _, model := range _list {
		list = append(list, ListResponse{
			NetModel:   model,
			NodeTitle:  model.NodeModel.Title,
			NodeStatus: model.NodeModel.Status,
		})
	}
	res.OkWithList(list, count, c)
}
