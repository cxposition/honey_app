package routers

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/internal/api"
	user_api2 "honey_app/apps/honey_server/internal/api/user_api"
	middleware2 "honey_app/apps/honey_server/internal/middleware"
)

func UserRouters(r *gin.RouterGroup) {
	var app = api.App.UserApi
	r.POST("/login", middleware2.BindJsonMiddleware[user_api2.LoginRequest], app.LoginView)
	r.POST("/users", middleware2.AdminMiddleware, middleware2.BindJsonMiddleware[user_api2.CreateRequest], app.CreateView)
	r.GET("/users", middleware2.BindQueryMiddleware[user_api2.UserListRequest], app.UserlistView)
	r.POST("/logout", app.UserLogoutView)
	r.DELETE("/users", middleware2.BindJsonMiddleware[user_api2.UserRemoveRequest], app.UserRemoveView)
	r.GET("/users/info", app.UserInfoView)
}
