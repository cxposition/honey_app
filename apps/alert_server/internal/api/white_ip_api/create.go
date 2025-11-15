package white_ip_api

import (
	"alert_server/internal/global"
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"alert_server/internal/utils/res"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	IP     string `json:"ip" binding:"required,ip"`
	Notice string `json:"notice"`
}

func (WhiteIPApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateRequest](c)
	var model models.WhiteIPModel
	err := global.DB.Take(&model, "ip = ?", cr.IP).Error
	if err == nil {
		res.FailWithMsg("白名单ip不能重复", c)
		return
	}
	err = global.DB.Create(&models.WhiteIPModel{
		IP:     cr.IP,
		Notice: cr.Notice,
	}).Error
	if err != nil {
		res.FailWithMsg("白名单ip保存失败", c)
		return
	}

	res.OkWithMsg("白名单ip保存成功", c)
}
