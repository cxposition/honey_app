package white_ip_api

import (
	"alert_server/internal/global"
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"alert_server/internal/utils/res"
	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	ID     uint   `json:"id" binding:"required"`
	IP     string `json:"ip" binding:"required"`
	Notice string `json:"notice"`
}

func (WhiteIPApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)
	var model models.WhiteIPModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("白名单ip不存在", c)
		return
	}

	var newModel models.WhiteIPModel
	err = global.DB.Take(&newModel, "id <> ? and ip = ?", cr.ID, cr.IP).Error
	if err == nil {
		res.OkWithData("修改的ip不能重复", c)
		return
	}

	err = global.DB.Model(&model).Updates(map[string]any{
		"ip":     cr.IP,
		"notice": cr.Notice,
	}).Error
	if err != nil {
		res.FailWithMsg("白名单更新失败", c)
		return
	}

	res.OkWithMsg("白名单更新成功", c)
}
