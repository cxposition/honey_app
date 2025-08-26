package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/utils"
	"honey_app/apps/honey_server/utils/jwts"
	"honey_app/apps/honey_server/utils/res"
)

func AuthMiddleware(c *gin.Context) {
	// 判断路径在不在白名单中，白名单中的路由不需要权限验证
	if utils.Inlist(global.Config.WhiteList, c.Request.URL.Path) {
		// 在白名单中直接放行
		c.Next()
		fmt.Println("免认证")
		return
	}
	token := c.GetHeader("token")
	_, err := jwts.ParseToken(token)
	if err != nil {
		logrus.Errorf("认证失败: %v", err)
		res.FailWithMsg("解析token失败", c)
		c.Abort() // 拦截
		return
	}
	c.Next()
	fmt.Println("认证通过")
}

func AdminMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ParseToken(token)
	if err != nil {
		c.JSON(200, gin.H{"code": 7, "msg": "认证失败", "data": gin.H{}})
		c.Abort()
		return
	}
	if claims.Role != 1 {
		c.JSON(200, gin.H{"code": 7, "msg": "角色认证失败", "data": gin.H{}})
		c.Abort()
		return
	}
	c.Next()
}
