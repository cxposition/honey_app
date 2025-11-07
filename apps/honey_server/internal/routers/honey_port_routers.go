package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/honey_port_api"
	"honey_server/internal/middleware"
)

func HoneyPortRouters(r *gin.RouterGroup) {
	app := api.App.HoneyPortApi
	r.PUT("honey_port", middleware.BindJsonMiddleware[honey_port_api.UpdateRequest], app.UpdateView)
	r.GET("honey_port", middleware.BindQueryMiddleware[honey_port_api.ListRequest], app.ListView)
}
