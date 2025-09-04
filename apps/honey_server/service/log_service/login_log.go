package log_service

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/models"
)

type LoginLogService struct {
	IP   string
	Addr string
}

func NewSuccessLog(c *gin.Context) *LoginLogService {
	return &LoginLogService{
		IP:   c.ClientIP(),
		Addr: "",
	}
}

func (l *LoginLogService) SuccessLog(userID uint, username string) {
	l.save(userID, username, "", "登录成功", true)
}

func (l *LoginLogService) FailLog(username string, password string, title string) {
	l.save(0, username, password, title, false)
}

func (l *LoginLogService) save(userID uint, username string, password string, title string, loginStatus bool) {
	global.DB.Create(&models.LogModel{
		Type:        1,
		IP:          l.IP,
		Addr:        l.Addr,
		UserID:      userID,
		Username:    username,
		Pwd:         password,
		LoginStatus: loginStatus,
		Title:       title,
	})
}
