package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("honey_server")
	g.Use(middleware.LogMiddleware, middleware.AuthMiddleware)
	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
