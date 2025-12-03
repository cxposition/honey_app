package api

import (
	"alert_server/internal/api/alert_api"
	"alert_server/internal/api/index_api"
	"alert_server/internal/api/white_ip_api"
)

type Api struct {
	WhiteIPApi white_ip_api.WhiteIPApi
	AlertApi   alert_api.AlertApi
	IndexApi   index_api.IndexApi
}

var App = Api{}
