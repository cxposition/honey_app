package user_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

type UserListRequest struct {
	models.PageInfo
	Username string `form:"username"`
}

func (UserApi) UserlistView(c *gin.Context) {
	cr := middleware.GetBind[UserListRequest](c)
	list, count, _ := common_service.QueryList(models.UserModel{}, common_service.ListRequest{
		Debug:    true,
		Likes:    []string{"username"}, // username like req.Key
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
