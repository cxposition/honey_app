package matrix_template_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"
)

type ListResponse struct {
	models.MatrixTemplateModel
	HostTemplateList []HostTemplateInfo `json:"hostTemplatePortInfo"`
}

type HostTemplateInfo struct {
	HostTemplateID    uint   `json:"hostTemplateID"`
	HostTemplateTitle string `json:"hostTemplateTitle"`
	Weight            int    `json:"weight"`
}

func (MatrixTemplateApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[models.PageInfo](c)
	_list, count, _ := common_service.QueryList(models.MatrixTemplateModel{},
		common_service.RequestList{
			Likes:    []string{"title"},
			Sort:     "created_at desc",
			PageInfo: cr,
		},
	)
	var list = make([]ListResponse, 0)
	var hostTemps []models.HostTemplateModel
	var hostTempIDList []uint
	for _, model := range _list {
		for _, port := range model.HostTemplateList {
			hostTempIDList = append(hostTempIDList, port.HostTemplateID)
		}
	}
	global.DB.Find(&hostTemps, "id in ?", hostTempIDList)
	var hostTempMap = map[uint]models.HostTemplateModel{}
	for _, i2 := range hostTemps {
		hostTempMap[i2.ID] = i2
	}
	for _, model := range _list {
		hostTemplateList := make([]HostTemplateInfo, 0)
		for _, port := range model.HostTemplateList {
			hostTemplateList = append(hostTemplateList, HostTemplateInfo{
				HostTemplateID:    port.HostTemplateID,
				HostTemplateTitle: hostTempMap[port.HostTemplateID].Title,
				Weight:            port.Weight,
			})
		}
		list = append(list, ListResponse{
			MatrixTemplateModel: model,
			HostTemplateList:    hostTemplateList,
		})
	}
	res.OkWithList(list, count, c)
}
