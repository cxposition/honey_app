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
	fmt.Println("offset:", offset, "cr.limit:", cr.Limit, "cr.Page", cr.Page)
	global.DB.Offset(offset).Limit(cr.Limit).Find(&list)
	var count int64
	global.DB.Model(&models.UserModel{}).Count(&count)
	res.OkWithList(list, count, c)
}
