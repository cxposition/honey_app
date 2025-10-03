package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/host_template_api"
	"image_server/internal/middleware"
	"image_server/internal/models"
)

func HostTemplateRouter(r *gin.RouterGroup) {
	app := api.App.HostTemplateApi
	r.POST("host_template", middleware.BindJsonMiddleware[host_template_api.CreateReuqest], app.CreateView)
	r.GET("host_template", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)
}
