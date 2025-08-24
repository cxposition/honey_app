package api

import "honey_app/apps/honey_server/api/user_api"

type Api struct {
	UserApi user_api.UserApi
}

var App = Api{}
