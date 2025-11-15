package white_ip_api

import (
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"alert_server/internal/service/common_service"
	"alert_server/internal/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (WhiteIPApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequestList](c)
	successCount, err := common_service.Remove(&models.WhiteIPModel{}, common_service.RemoveRequest{
		IDList: cr.IdList,
		Log:    middleware.GetLog(c),
		Msg:    "白名单",
	})
	if err != nil {
		res.FailWithMsg(fmt.Sprintf("白名单删除失败 %s", err), c)
		return
	}

	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IdList), successCount)
	res.OkWithMsg(msg, c)
	return
}
