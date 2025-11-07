package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/honey_ip_api"
	"honey_server/internal/middleware"
)

func HoneyIPRouters(r *gin.RouterGroup) {
	app := api.App.HoneyIPApi
	r.POST("honey_ip", middleware.BindJsonMiddleware[honey_ip_api.CreateRequest], app.CreateView)
	r.GET("honey_ip", middleware.BindQueryMiddleware[honey_ip_api.ListRequest], app.ListView)
}
