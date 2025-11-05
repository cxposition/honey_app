package host_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

func (HostApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequestList](c)
	log := middleware.GetLog(c)
	successCount, err := common_service.Remove(models.HostModel{}, common_service.RemoveRequest{
		IDList: cr.IdList,
		Log:    log,
		Msg:    "主机",
	})
	if err != nil {
		msg := fmt.Sprintf("网络删除失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}
	msg := fmt.Sprintf("成功删除 %d 个存活主机", successCount)
	res.OkWithMsg(msg, c)
}
