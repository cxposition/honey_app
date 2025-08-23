package routers

import "github.com/gin-gonic/gin"

func UserRouters(r *gin.RouterGroup) {
	r.GET("users", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0})
	})
	r.GET("login", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 1})
	})
}
