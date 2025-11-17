package routers

import (
	"alert_server/internal/api"
	"alert_server/internal/api/alert_api"
	"alert_server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AlertRouter(r *gin.RouterGroup) {
	app := api.App.AlertApi
	r.GET("alert", middleware.BindQueryMiddleware[alert_api.ListRequest], app.ListView)
}
