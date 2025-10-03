package host_template_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type OptionsListResponse struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

func (HostTemplateApi) OptionsView(c *gin.Context) {
	var list = make([]OptionsListResponse, 0)
	global.DB.Model(models.HostTemplateModel{}).Select("id as value", "title as label").Scan(&list)
	res.OkWithData(list, c)
}
