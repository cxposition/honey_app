package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/net_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
)

func NetRouter(r *gin.RouterGroup) {
	var app = api.App.NetApi
	r.GET("net", middleware.BindQueryMiddleware[net_api.ListRequest], app.ListView)
	r.GET("net/options", app.OptionsView)
	r.GET("net/:id", middleware.BindUriMiddleware[models.IDRequest], app.DetailView)
	r.PUT("net", middleware.BindJsonMiddleware[net_api.UpdateRequest], app.UpdateView)
	r.DELETE("net", middleware.BindJsonMiddleware[models.IDRequestList], app.RemoveView)
	r.POST("net/scan", middleware.BindJsonMiddleware[models.IDRequest], app.ScanView)
	r.GET("net/ip_list", middleware.BindQueryMiddleware[net_api.NetUseIPListRequest], app.NetUseIPListView)
}
