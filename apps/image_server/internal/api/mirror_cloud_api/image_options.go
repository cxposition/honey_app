package mirror_cloud_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type ImageOptionsListResponse struct {
	Label   string `json:"label"`
	Value   uint   `json:"value"`
	Disable bool   `json:"disable"`
}

func (MirrorCloudApi) ImageOptionsListView(c *gin.Context) {
	var list []models.ImageModel
	global.DB.Unscoped().Find(&list)
	var options []ImageOptionsListResponse
	fmt.Println("list:", list)
	for _, model := range list {
		item := ImageOptionsListResponse{
			Label:   fmt.Sprintf("%s/%d", model.Title, model.Port),
			Value:   model.ID,
			Disable: false,
		}
		if model.Status == 2 {
			item.Disable = true
		}
		options = append(options, item)
	}
	res.OkWithData(options, c)
}
