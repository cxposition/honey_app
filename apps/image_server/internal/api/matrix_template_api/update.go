package matrix_template_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type UpdateRequest struct {
	ID       uint                        `json:"id" binding:"required"`
	Title    string                      `json:"title" binding:"required"`
	PortList models.HostTemplatePortList `json:"portList" binding:"dive"`
}

func (MatrixTemplateApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)
	var model models.HostTemplateModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("主机模版不存在", c)
		return
	}

	var newModel models.HostTemplateModel
	err = global.DB.Take(&newModel, "title = ? and id <> ?", cr.Title, cr.ID).Error
	if err == nil {
		res.FailWithMsg("修改主机名称不能重复", c)
		return
	}

	// 校验服务id
	// 校验端口不能重复
	var serviceIDList []uint
	var portMap = map[int]bool{}
	for _, port := range cr.PortList {
		serviceIDList = append(serviceIDList, port.ServiceID)
		portMap[port.Port] = true
	}

	if len(portMap) != len(cr.PortList) {
		res.FailWithMsg("端口存在重复", c)
		return
	}

	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	var serviceMap = map[uint]models.ServiceModel{}
	for _, serviceModel := range serviceList {
		serviceMap[serviceModel.ID] = serviceModel
	}
	for _, port := range cr.PortList {
		model, ok := serviceMap[port.ServiceID]
		if !ok {
			msg := fmt.Sprintf("服务id不存在: %d", port.ServiceID)
			res.FailWithMsg(msg, c)
			return
		}
		if model.Status != 1 {
			res.FailWithMsg("服务状态异常", c)
			return
		}
	}

	// 消息入库
	newModel = models.HostTemplateModel{
		Title:    cr.Title,
		PortList: cr.PortList,
	}
	err = global.DB.Model(&model).Updates(map[string]any{
		"title":     cr.Title,
		"port_list": cr.PortList,
	}).Error
	if err != nil {
		res.FailWithMsg("主机模版更新失败", c)
		return
	}
	res.OkWithMsg("主机模版更新成功", c)
}
