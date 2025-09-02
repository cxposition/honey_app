package user_api

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/utils/res"
	"time"
)

func (UserApi) UserLogoutView(c *gin.Context) {
	token := c.GetHeader("token")
	log := middleware.GetLog(c)
	auth := middleware.GetAuth(c)
	expiresAt := time.Unix(auth.ExpiresAt, 0)
	log.Infof("用户注销 %d %s %s", auth.UserID, token, expiresAt)
	res.OkWithMsg("注销成功", c)
}
