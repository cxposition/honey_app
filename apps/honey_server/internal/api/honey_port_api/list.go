package honey_port_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

type ListRequest struct {
	models.PageInfo
	HoneyIPID uint `form:"honeyIpID" binding:"required"`
}

type ListResponse struct {
	models.HoneyPortModel
	ServiceTitle string `json:"serviceTitle"`
}

func (HoneyPortApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
	_list, count, _ := common_service.QueryList(models.HoneyPortModel{HoneyIpID: cr.HoneyIPID}, common_service.ListRequest{
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
		Preload:  []string{"ServiceModel"},
	})
	var list = make([]ListResponse, 0)
	for _, model := range _list {
		list = append(list, ListResponse{
			HoneyPortModel: model,
			ServiceTitle:   model.ServiceModel.Title,
		})
	}
	res.OkWithList(list, count, c)
}
