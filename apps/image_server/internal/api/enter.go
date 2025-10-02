package api

import (
	"image_server/internal/api/host_template_api"
	"image_server/internal/api/mirror_cloud_api"
	"image_server/internal/api/vs_api"
	"image_server/internal/api/vs_net_api"
)

type Api struct {
	MirrorCloudApi  mirror_cloud_api.MirrorCloudApi
	VsApi           vs_api.VsApi
	VsNetApi        vs_net_api.VsNetApi
	HostTemplateApi host_template_api.HostTemplateApi
}

var App = Api{}
