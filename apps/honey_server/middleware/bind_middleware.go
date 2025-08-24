package middleware

import "github.com/gin-gonic/gin"

func BindJsonMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": "参数绑定错误"})
		c.Abort()
		return
	}
	c.Set("request", cr)
}

func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		c.JSON(200, gin.H{"code": 1, "msg": "参数绑定错误"})
		c.Abort()
		return
	}
	c.Set("request", cr)
}

func GetBind[T any](c *gin.Context) T {
	return c.MustGet("request").(T)
}
