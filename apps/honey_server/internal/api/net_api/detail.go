package net_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

func (NetApi) DetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NetModel
	logrus.Infof("获取网络详情,ID:%d", cr.ID)
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}
	res.OkWithData(model, c)
}
