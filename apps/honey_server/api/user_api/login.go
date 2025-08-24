package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/middleware"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (UserApi) LoginView(c *gin.Context) {
	cr := middleware.GetBind[LoginRequest](c)
	fmt.Printf("cr: %+v\n", cr)
	c.JSON(200, gin.H{"code": 0, "msg": "登陆成功"})
}
