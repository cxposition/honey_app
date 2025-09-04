package api

import (
	"honey_app/apps/honey_server/api/captcha_api"
	"honey_app/apps/honey_server/api/log_api"
	"honey_app/apps/honey_server/api/user_api"
)

type Api struct {
	UserApi    user_api.UserApi
	CaptchaApi captcha_api.CaptchaApi
	LogApi     log_api.LogApi
}

var App = Api{}
