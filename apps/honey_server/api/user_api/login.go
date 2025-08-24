package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/utils/res"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (UserApi) LoginView(c *gin.Context) {
	cr := middleware.GetBind[LoginRequest](c)
	fmt.Printf("cr: %+v\n", cr)
	res.OkWithMsg("登录成功", c)
}
