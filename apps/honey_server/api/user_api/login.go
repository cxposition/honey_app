package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/utils/captcha"
	"honey_app/apps/honey_server/utils/jwts"
	"honey_app/apps/honey_server/utils/res"
)

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captchaId" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

func (UserApi) LoginView(c *gin.Context) {
	cr := middleware.GetBind[LoginRequest](c)
	if cr.CaptchaID == "" || cr.CaptchaCode == "" {
		res.FailWithMsg("请输入图片验证码", c)
		return
	}
	if !captcha.CaptchaStore.Verify(cr.CaptchaID, cr.CaptchaCode, true) {
		res.FailWithMsg("图片验证码验证失败", c)
		return
	}
	var user models.UserModel
	err := global.DB.Take(&user, "username = ?", cr.Username).Error
	if err != nil {
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	token, err := jwts.GetToken(jwts.ClaimsUserInfo{
		UserID: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		logrus.Errorf("生成token失败 %s", err)
		res.FailWithMsg("登陆失败", c)
		return
	}

	res.OkWithData(token, c)
}
