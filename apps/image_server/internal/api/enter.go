package api

import (
	"image_server/internal/api/mirror_cloud_api"
	"image_server/internal/api/vs_api"
)

type Api struct {
	MirrorCloudApi mirror_cloud_api.MirrorCloudApi
	VsApi          vs_api.VsApi
}

var App = Api{}
