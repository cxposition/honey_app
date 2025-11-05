package host_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

type ListRequest struct {
	models.PageInfo
	NodeID uint `form:"nodeID"`
	NetID  uint `form:"netID"`
}

type ListResponse struct {
	models.HostModel
	NetTitle  string `json:"netTitle"`
	NodeTitle string `json:"nodeTitle"`
}

func (HostApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	_list, count, _ := common_service.QueryList(models.HostModel{NodeID: cr.NodeID, NetID: cr.NetID}, common_service.Request{
		Likes:    []string{"ip", "mac"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
		Preload:  []string{"NodeModel", "NetModel"},
	})
	var list = make([]ListResponse, 0, len(_list))
	for _, model := range _list {
		list = append(list, ListResponse{
			HostModel: model,
			NodeTitle: model.NodeModel.Title,
			NetTitle:  model.NetModel.Title,
		})
	}
	res.OkWithList(list, count, c)
}
