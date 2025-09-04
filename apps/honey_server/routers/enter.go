package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/middleware"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("honey_server")
	g.Use(middleware.LogMiddleware, middleware.AuthMiddleware)
	UserRouters(g)
	CaptchaRouters(g)
	LogRouters(g)
	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
