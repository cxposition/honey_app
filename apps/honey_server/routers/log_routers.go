package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/api"
	"honey_app/apps/honey_server/api/log_api"
	"honey_app/apps/honey_server/middleware"
)

func LogRouters(r *gin.RouterGroup) {
	var app = api.App.LogApi
	r.GET("logs", middleware.AdminMiddleware, middleware.BindQueryMiddleware[log_api.LogListRequest], app.LogListView)
}
