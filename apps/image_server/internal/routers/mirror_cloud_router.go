package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
)

func MirrorCloudRouter(r *gin.RouterGroup) {
	app := api.App.MirrorCloudApi
	r.POST("mirror_cloud/see", app.ImageSeeView)
}
