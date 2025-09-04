package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/middleware"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/service/log_service"
	"honey_app/apps/honey_server/utils/captcha"
	"honey_app/apps/honey_server/utils/jwts"
	"honey_app/apps/honey_server/utils/pwd"
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
	LoginLog := log_service.NewSuccessLog(c)
	if cr.CaptchaID == "" || cr.CaptchaCode == "" {
		LoginLog.FailLog(cr.Username, "", "未输入图片验证码")
		res.FailWithMsg("请输入图片验证码", c)
		return
	}
	if !captcha.CaptchaStore.Verify(cr.CaptchaID, cr.CaptchaCode, true) {
		LoginLog.FailLog(cr.Username, "", "图片验证码验证失败")
		res.FailWithMsg("图片验证码验证失败", c)
		return
	}
	var user models.UserModel
	err := global.DB.Take(&user, "username = ?", cr.Username).Error
	if err != nil {
		LoginLog.FailLog(cr.Username, cr.Password, "用户名不存在")
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		LoginLog.FailLog(cr.Username, cr.Password, "密码错误")
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

	LoginLog.SuccessLog(user.ID, user.Username)
	res.OkWithData(token, c)
}
