package log_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/service/common_service"
	"honey_app/apps/honey_server/utils/res"
)

type LogRemoveRequest struct {
	IDList []uint `json:"idList"`
}

func (l *LogApi) LogRemoveView(c *gin.Context) {
	cr := middleware.GetBind[LogRemoveRequest](c)
	log := middleware.GetLog(c)
	successCount, err := common_service.Remove(models.LogModel{}, common_service.RemoveRequest{
		Debug:    true,
		IDList:   cr.IDList,
		Log:      log,
		Msg:      "登陆日志",
		Unscoped: true,
	})
	if err != nil {
		msg := fmt.Sprintf("删除登陆日志失败: %s", err)
		res.FailWithMsg(msg, c)
		return
	}
	msg := fmt.Sprintf("删除成功, 共%d个, 成功 %d 个", len(cr.IDList), successCount)
	res.OkWithMsg(msg, c)
}
