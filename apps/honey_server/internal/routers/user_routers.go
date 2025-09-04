package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/api/user_api"
	"honey_server/internal/middleware"
)

func UserRouters(r *gin.RouterGroup) {
	var app = api.App.UserApi
	r.POST("/login", middleware.BindJsonMiddleware[user_api.LoginRequest], app.LoginView)
	r.POST("/users", middleware.AdminMiddleware, middleware.BindJsonMiddleware[user_api.CreateRequest], app.CreateView)
	r.GET("/users", middleware.BindQueryMiddleware[user_api.UserListRequest], app.UserlistView)
	r.POST("/logout", app.UserLogoutView)
	r.DELETE("/users", middleware.BindJsonMiddleware[user_api.UserRemoveRequest], app.UserRemoveView)
	r.GET("/users/info", app.UserInfoView)
}
