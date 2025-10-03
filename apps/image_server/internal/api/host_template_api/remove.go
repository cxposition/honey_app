package host_template_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"
)

func (HostTemplateApi) Remove(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)
	log := middleware.GetLog(c)
	successCount, err := common_service.Remove(models.HostTemplateModel{}, common_service.RemoveRequest{
		Debug:  true,
		IDList: cr.IDList,
		Log:    log,
		Msg:    "主机模板",
	})
	if err != nil {
		msg := fmt.Sprintf("删除主机模板失败: %s", err)
		res.FailWithMsg(msg, c)
		return
	}
	msg := fmt.Sprintf("删除成功, 共%d个, 成功 %d 个", len(cr.IDList), successCount)
	res.OkWithMsg(msg, c)
}
