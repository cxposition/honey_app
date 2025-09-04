package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/log_api"
	"honey_server/internal/middleware"
)

func LogRouters(r *gin.RouterGroup) {
	var app = api.App.LogApi
	r.GET("logs", middleware.AdminMiddleware, middleware.BindQueryMiddleware[log_api.LogListRequest], app.LogListView)
	r.DELETE("logs", middleware.AdminMiddleware, middleware.BindJsonMiddleware[log_api.LogRemoveRequest], app.LogRemoveView)
}
