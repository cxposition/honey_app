package api

import (
	"honey_app/apps/honey_server/api/captcha_api"
	"honey_app/apps/honey_server/api/user_api"
)

type Api struct {
	UserApi    user_api.UserApi
	CaptchaApi captcha_api.CaptchaApi
}

var App = Api{}
