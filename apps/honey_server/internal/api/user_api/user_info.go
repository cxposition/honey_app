package user_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

type UserInfoResponse struct {
	UserID        uint   `json:"userID"`
	Username      string `json:"username"`
	Role          int8   `json:"role"` // 1 管理员 2 普通用户
	LastLoginDate string `json:"lastLoginDate"`
}

func (UserApi) UserInfoView(c *gin.Context) {
	auth := middleware.GetAuth(c)
	var user models.UserModel
	err := global.DB.Take(&user, auth.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	data := UserInfoResponse{
		UserID:        user.ID,
		Username:      user.Username,
		Role:          user.Role,
		LastLoginDate: user.LastLoginDate,
	}
	res.OkWithData(data, c)
}
