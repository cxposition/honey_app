package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/global"
)

func Run() {
	r := gin.Default()
	r.Static("uploads", "uploads")
	g := r.Group("honey_server")
	g.Use()

	UserRouters(g)
	webAddr := global.Config.System.WebAddr
	r.Run(webAddr)
}
