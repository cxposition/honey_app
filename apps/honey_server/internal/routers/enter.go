package routers

import (
	ginpprof "github.com/gin-contrib/pprof"
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
	HoneyPortRouters(g)

	webAddr := global.Config.System.WebAddr
	go ginpprof.Register(r)
	r.Run(webAddr)
}
