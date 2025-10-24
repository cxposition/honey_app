package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
)

func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi
	r.GET("node_network/flush", app.FlushView)
}
