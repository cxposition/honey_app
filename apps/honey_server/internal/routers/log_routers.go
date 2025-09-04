package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/internal/api"
	log_api2 "honey_app/apps/honey_server/internal/api/log_api"
	middleware2 "honey_app/apps/honey_server/internal/middleware"
)

func LogRouters(r *gin.RouterGroup) {
	var app = api.App.LogApi
	r.GET("logs", middleware2.AdminMiddleware, middleware2.BindQueryMiddleware[log_api2.LogListRequest], app.LogListView)
	r.DELETE("logs", middleware2.AdminMiddleware, middleware2.BindJsonMiddleware[log_api2.LogRemoveRequest], app.LogRemoveView)
}
