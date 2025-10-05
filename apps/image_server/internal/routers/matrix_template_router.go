package routers

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/api"
	"image_server/internal/api/matrix_template_api"
	"image_server/internal/middleware"
	"image_server/internal/models"
)

func MatrixTemplateRouter(r *gin.RouterGroup) {
	app := api.App.MatrixTemplateApi
	r.POST("matrix_template", middleware.BindJsonMiddleware[matrix_template_api.CreateReuqest], app.CreateView)
	r.PUT("matrix_template", middleware.BindJsonMiddleware[matrix_template_api.UpdateRequest], app.UpdateView)
	r.GET("matrix_template", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)
	r.GET("matrix_template/options", app.OptionsView)
	r.DELETE("matrix_template", middleware.BindJsonMiddleware[models.IDListRequest], app.Remove)
}
