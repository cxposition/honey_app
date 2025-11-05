package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/host_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
)

func HostRouters(r *gin.RouterGroup) {
	app := api.App.HostApi
	r.GET("host", middleware.BindQueryMiddleware[host_api.ListRequest], app.ListView)
	r.DELETE("host", middleware.BindJsonMiddleware[models.IDRequestList], app.RemoveView)
}
