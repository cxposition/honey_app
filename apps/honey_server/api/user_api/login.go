package user_api

import (
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
	log := middleware.GetLog(c)
	log.Infof("这是请求的内容:%v", cr)
	res.OkWithMsg("登录成功", c)
}
