package vs_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type VsOptionsListResponse struct {
	Label   string `json:"label"`
	Value   uint   `json:"value"`
	Disable bool   `json:"disable"`
}

func (VsApi) VsOptionsListView(c *gin.Context) {
	var list []models.ServiceModel
	global.DB.Unscoped().Find(&list)
	var options []VsOptionsListResponse
	fmt.Println("list:", list)
	for _, model := range list {
		item := VsOptionsListResponse{
			Label:   fmt.Sprintf("%s/%d", model.Title, model.Port),
			Value:   model.ID,
			Disable: false,
		}
		if model.Status != 1 {
			item.Disable = true
		}
		options = append(options, item)
	}
	res.OkWithData(options, c)
}
