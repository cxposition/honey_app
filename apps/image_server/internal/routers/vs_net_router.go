package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/vs_net_api"
	"image_server/internal/middleware"
)

func VsNetRouter(r *gin.RouterGroup) {
	app := api.App.VsNetApi
	r.PUT("vs_net", middleware.BindJsonMiddleware[vs_net_api.VsNetRequest], app.VsNetUpdateView)
	r.GET("vs_net", app.VsNetInfoView)
}
