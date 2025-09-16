package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/vs_api"
	"image_server/internal/middleware"
)

func VsRouter(r *gin.RouterGroup) {
	app := api.App.VsApi
	r.POST("vs", middleware.BindJsonMiddleware[vs_api.VsCreateRequest], app.VsCreateView)
	r.GET("vs", middleware.BindQueryMiddleware[vs_api.VsListRequest], app.VsListView)
}
