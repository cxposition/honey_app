package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
)

func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi
	r.GET("node_network/flush", middleware.BindQueryMiddleware[models.IDRequest], app.FlushView)
}
