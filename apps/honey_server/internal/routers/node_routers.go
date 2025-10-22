package routers

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
)

func NodeRouters(r *gin.RouterGroup) {
	var app = api.App.NodeApi
	r.GET("node", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)
}
