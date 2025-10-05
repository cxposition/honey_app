package matrix_template_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type CreateReuqest struct {
	Title            string                      `json:"title" binding:"required"`
	PortList         models.HostTemplatePortList `json:"portList"`
	HostTemplateList models.HostTemplateList     `gorm:"serializer:json" json:"HostTemplateList" binding:"required,dive"`
}

func (MatrixTemplateApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateReuqest](c)
	var model models.MatrixTemplateModel
	err := global.DB.Take(&model, "title = ?", cr.Title).Error
	if err == nil {
		res.FailWithMsg("矩阵模版名称不能重复", c)
		return
	}

	// 校验服务id
	// 校验端口不能重复
	var hostTemplateIDList []uint
	for _, h := range cr.HostTemplateList {
		hostTemplateIDList = append(hostTemplateIDList, h.HostTemplateID)
	}

	var hostTemps []models.HostTemplateModel
	global.DB.Find(&hostTemps, "id in ?", hostTemplateIDList)
	var hostTempMap = map[uint]models.HostTemplateModel{}
	for _, m := range hostTemps {
		hostTempMap[m.ID] = m
	}
	for _, h := range cr.HostTemplateList {
		_, ok := hostTempMap[h.HostTemplateID]
		if !ok {
			msg := fmt.Sprintf("服务id不存在: %d", h.HostTemplateID)
			res.FailWithMsg(msg, c)
			return
		}
	}

	// 消息入库
	model = models.MatrixTemplateModel{
		Title:            cr.Title,
		HostTemplateList: cr.HostTemplateList,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("矩阵模版创建失败", c)
		return
	}
	res.OkWithData(model.ID, c)
}
