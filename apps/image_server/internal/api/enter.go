package api

import "image_server/internal/api/mirror_cloud_api"

type Api struct {
	MirrorCloudApi mirror_cloud_api.MirrorCloudApi
}

var App = Api{}
