package middleware

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/utils/res"
)

func BindJsonMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
}

func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithMsg("参数绑定错误", c)
		c.Abort()
		return
	}
	c.Set("request", cr)
}

func BindUriMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBind(&cr)
	if err != nil {
		res.FailWithMsg("参数绑定错误", c)
		c.Abort()
		return
	}
	c.Set("request", cr)
}

func GetBind[T any](c *gin.Context) T {
	return c.MustGet("request").(T)
}
