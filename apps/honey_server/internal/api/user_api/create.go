package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/service/user_service"
	"honey_server/internal/utils/res"
)

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     int8   `json:"role" binding:"required,ne=1"`
}

func (UserApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateRequest](c)
	log := middleware.GetLog(c)
	us := user_service.NewUserService(log)
	user, err := us.Create(user_service.UserCreateRequest{
		Password: cr.Password,
		Role:     cr.Role,
		Username: cr.Username,
	})
	if err != nil {
		msg := fmt.Sprintf("创建用户失败 %s", err)
		log.Errorf(msg)
		res.FailWithMsg(msg, c)
	}
	res.OkWithData(user, c)
}
