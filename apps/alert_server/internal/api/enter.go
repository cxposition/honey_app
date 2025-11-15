package api

import "alert_server/internal/api/white_ip_api"

type Api struct {
	WhiteIPApi white_ip_api.WhiteIPApi
}

var App = Api{}
