package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/node_network_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
)

func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi
	r.GET("node_network/flush", middleware.BindQueryMiddleware[models.IDRequest], app.FlushView)
	r.GET("node_network", middleware.BindQueryMiddleware[node_network_api.ListRequest], app.ListView)
	r.PUT("node_network", middleware.BindJsonMiddleware[node_network_api.UpdateRequest], app.UpdateView)
	r.PUT("node_network/enable", middleware.BindJsonMiddleware[models.IDRequest], app.EnableView)
}
