package captcha_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/utils/captcha"
	"honey_app/apps/honey_server/utils/res"
)

type CaptchaApi struct {
}

type GenerateResponse struct {
	CaptchaID string `json:"captchaID"`
	Captcha   string `json:"captcha"`
}

func (CaptchaApi) GenerateView(c *gin.Context) {
	var driver = base64Captcha.DriverString{
		Width:           200,
		Height:          60,
		NoiseCount:      2,
		ShowLineOptions: 4,
		Length:          4,
		Source:          "0123456789",
	}
	cp := base64Captcha.NewCaptcha(&driver, captcha.CaptchaStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		logrus.Errorf("图片验证码生成失败 %s", err)
		res.FailWithMsg("图片验证码生成失败", c)
		return
	}
	res.OkWithData(GenerateResponse{
		CaptchaID: id,
		Captcha:   b64s,
	}, c)
}
