package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/mirror_cloud_api"
	"image_server/internal/middleware"
)

func MirrorCloudRouter(r *gin.RouterGroup) {
	app := api.App.MirrorCloudApi
	r.POST("mirror_cloud/see", app.ImageSeeView)
	r.POST("mirror_cloud", middleware.BindJsonMiddleware[mirror_cloud_api.ImageCreateRequest], app.ImageCreateView)
	r.GET("mirror_cloud", middleware.BindQueryMiddleware[mirror_cloud_api.ImageListRequest], app.ImageListView)
}
