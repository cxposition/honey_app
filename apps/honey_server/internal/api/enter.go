package api

import (
	"honey_server/internal/api/captcha_api"
	"honey_server/internal/api/host_api"
	"honey_server/internal/api/log_api"
	"honey_server/internal/api/net_api"
	"honey_server/internal/api/node_api"
	"honey_server/internal/api/node_network_api"
	"honey_server/internal/api/user_api"
)

type Api struct {
	UserApi        user_api.UserApi
	CaptchaApi     captcha_api.CaptchaApi
	LogApi         log_api.LogApi
	NodeApi        node_api.NodeApi
	NodeNetworkApi node_network_api.NodeNetworkApi
	NetApi         net_api.NetApi
	HostApi        host_api.HostApi
}

var App = Api{}
