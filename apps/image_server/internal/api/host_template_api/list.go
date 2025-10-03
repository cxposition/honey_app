package host_template_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"
)

type ListResponse struct {
	models.HostTemplateModel
	PortList []HostTemplatePortInfo `json:"portList"`
}

type HostTemplatePortInfo struct {
	Port          int    `json:"port"`
	ServiceID     uint   `json:"serviceID"`
	ServiceTitle  string `json:"serviceTitle"`
	ServiceStatus int8   `json:"serviceStatus"`
}

func (HostTemplateApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[models.PageInfo](c)
	_list, count, _ := common_service.QueryList(models.HostTemplateModel{},
		common_service.RequestList{
			Likes:    []string{"title"},
			Sort:     "created_at desc",
			PageInfo: cr,
		},
	)
	var list = make([]ListResponse, 0)
	var serviceList []models.ServiceModel
	var serviceIDList []uint
	for _, model := range _list {
		for _, port := range model.PortList {
			serviceIDList = append(serviceIDList, port.ServiceID)
		}
	}
	fmt.Println(serviceIDList)
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	fmt.Println(serviceList)
	var serviceMap = map[uint]models.ServiceModel{}
	for _, model := range _list {
		portList := make([]HostTemplatePortInfo, 0)
		for _, port := range model.PortList {
			portList = append(portList, HostTemplatePortInfo{
				Port:          port.Port,
				ServiceID:     port.ServiceID,
				ServiceTitle:  serviceMap[port.ServiceID].Title,
				ServiceStatus: serviceMap[port.ServiceID].Status,
			})
		}
		list = append(list, ListResponse{
			HostTemplateModel: model,
			PortList:          portList,
		})
	}
	res.OkWithList(list, count, c)
}
