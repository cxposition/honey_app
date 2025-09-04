package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware2 "honey_app/apps/honey_server/internal/middleware"
	"honey_app/apps/honey_server/internal/models"
	"honey_app/apps/honey_server/internal/service/common_service"
	"honey_app/apps/honey_server/internal/utils/res"
)

type UserRemoveRequest struct {
	IDList []uint `json:"idList"`
}

func (UserApi) UserRemoveView(c *gin.Context) {
	cr := middleware2.GetBind[UserRemoveRequest](c)
	log := middleware2.GetLog(c)
	successCount, err := common_service.Remove(models.UserModel{}, common_service.RemoveRequest{
		Debug:  true,
		IDList: cr.IDList,
		Log:    log,
		Msg:    "用户",
	})
	if err != nil {
		msg := fmt.Sprintf("删除用户失败: %s", err)
		res.FailWithMsg(msg, c)
		return
	}
	msg := fmt.Sprintf("删除成功, 共%d个, 成功 %d 个", len(cr.IDList), successCount)
	res.OkWithMsg(msg, c)
}
