package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/api"
	"honey_app/apps/honey_server/api/user_api"
	"honey_app/apps/honey_server/middleware"
)

func UserRouters(r *gin.RouterGroup) {
	var app = api.App.UserApi
	r.POST("/login", middleware.BindJsonMiddleware[user_api.LoginRequest], app.LoginView)
	r.POST("/users", middleware.AdminMiddleware, middleware.BindJsonMiddleware[user_api.CreateRequest], app.CreateView)
}
