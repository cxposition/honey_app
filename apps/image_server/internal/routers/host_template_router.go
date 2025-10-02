package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/host_template_api"
	"image_server/internal/middleware"
)

func HostTemplateRouter(r *gin.RouterGroup) {
	app := api.App.HostTemplateApi
	r.POST("host_template", middleware.BindJsonMiddleware[host_template_api.CreateReuqest], app.CreateView)
}
