package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/api"
)

func CaptchaRouters(r *gin.RouterGroup) {
	var app = api.App.CaptchaApi
	r.GET("captcha", app.GenerateView)
}
