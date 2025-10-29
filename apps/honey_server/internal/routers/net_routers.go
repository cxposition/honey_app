package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/net_api"
	"honey_server/internal/middleware"
)

func NetRouter(r *gin.RouterGroup) {
	var app = api.App.NetApi
	r.GET("net", middleware.BindQueryMiddleware[net_api.ListRequest], app.ListView)
}
