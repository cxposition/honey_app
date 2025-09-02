package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/utils/res"
)

type UserListRequest struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Key      string `form:"key"`
	Username string `form:"username"`
}

func (UserApi) UserlistView(c *gin.Context) {
	cr := middleware.GetBind[UserListRequest](c)
	fmt.Printf("cr:%+v", cr)
	list := make([]models.UserModel, 0)
	if cr.Limit <= 0 {
		cr.Limit = 10
	}
	if cr.Page <= 0 {
		cr.Page = 1
	}
	offset := (cr.Page - 1) * cr.Limit
	query := global.DB.Where(models.UserModel{
		Username: cr.Username,
	})
	var like = global.DB.Where("")
	if cr.Key != "" {
		like.Where("username like ?", fmt.Sprintf("%%%s%%", cr.Key))
	}
	baseDB := global.DB.Debug()
	baseDB = baseDB.Preload("LogList")

	baseDB.Offset(offset).Limit(cr.Limit).Where(like).Where(query).Find(&list)
	var count int64
	global.DB.Debug().Model(&models.UserModel{}).Where(like).Where(query).Count(&count)
	res.OkWithList(list, count, c)
}
