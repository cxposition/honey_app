package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("honey_server")
	g.Use(middleware.LogMiddleware, middleware.AuthMiddleware)

	UserRouters(g)
	CaptchaRouters(g)
	LogRouters(g)
	NodeRouters(g)
	NodeNetworkRouters(g)
	NetRouter(g)
	HostRouters(g)
	HoneyIPRouters(g)

	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
