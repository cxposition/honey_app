package user_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/utils/res"
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
