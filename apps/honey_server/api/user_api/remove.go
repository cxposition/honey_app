package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/service/common_service"
	"honey_app/apps/honey_server/utils/res"
)

type UserRemoveRequest struct {
	IDList []uint `json:"idList"`
}

func (UserApi) UserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[UserRemoveRequest](c)
	log := middleware.GetLog(c)
	//var list []models.UserModel
	//global.DB.Find(&list, "id in ?", cr.IDList)
	//var successCount = int64(0)
	//if len(list) > 0 {
	//	db := global.DB.Delete(&list)
	//	err := db.Error
	//	if err != nil {
	//		msg := fmt.Sprintf("删除用户失败: %s", db.Error)
	//		log.Errorf(msg)
	//		log.Errorf("删除失败的入参 %v", cr.IDList)
	//		res.FailWithMsg(msg, c)
	//		return
	//	}
	//	successCount = db.RowsAffected
	//}
	//
	//msg := fmt.Sprintf("删除用户成功共%d个，成功%d个", successCount, len(list))
	//log.Infof(msg)
	//log.Infof("删除的用户信息%v", list)
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
