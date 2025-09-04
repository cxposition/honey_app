package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	middleware2 "image_server/internal/middleware"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("honey_server")
	g.Use(middleware2.LogMiddleware, middleware2.AuthMiddleware)
	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
