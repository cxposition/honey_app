package routers

import (
	"alert_server/internal/api"
	"alert_server/internal/api/white_ip_api"
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"github.com/gin-gonic/gin"
)

func WhiteIPRouter(r *gin.RouterGroup) {
	app := api.Api{}.WhiteIPApi
	r.GET("white_ip", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)
	r.POST("white_ip", middleware.BindJsonMiddleware[white_ip_api.CreateRequest], app.CreateView)
	r.PUT("white_ip", middleware.BindJsonMiddleware[white_ip_api.UpdateRequest], app.UpdateView)
	r.DELETE("white_ip", middleware.BindJsonMiddleware[models.IDRequestList], app.RemoveView)
}
