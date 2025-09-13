package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
	"image_server/internal/middleware"
)

func Run() {
	system := global.Config.System
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("image_server")
	//g.Use(middleware.LogMiddleware, middleware.AuthMiddleware)
	g.Use(middleware.LogMiddleware)
	MirrorCloudRouter(g)
	VsRouter(g)
	webAddr := system.WebAddr
	logrus.Infof("web addr run %s", webAddr)
	r.Run(webAddr)
}
