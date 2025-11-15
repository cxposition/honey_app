package routers

import (
	"alert_server/internal/global"
	"alert_server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("alert_server")
	//g.Use(middleware.LogMiddleware, middleware.AuthMiddleware)
	g.Use(middleware.LogMiddleware)
	WhiteIPRouter(g)
	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
