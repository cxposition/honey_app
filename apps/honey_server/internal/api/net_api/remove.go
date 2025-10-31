package net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

func (NetApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequestList](c)
	log := middleware.GetLog(c)
	successCount, err := common_service.Remove(models.NetModel{}, common_service.RemoveRequest{
		IDList: cr.IdList,
		Log:    log,
		Msg:    "网络",
	})
	if err != nil {
		msg := fmt.Sprintf("网络删除失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}
	msg := fmt.Sprintf("成功删除 %d 个网络", successCount)
	res.OkWithMsg(msg, c)
}
