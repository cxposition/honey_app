package mirror_cloud_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type ImageUpdateRequest struct {
	ID        uint   `json:"id"`
	Title     string `json:"title" binding:"required"`
	Port      int    `json:"port" binding:"required,min=1,max=65535"`
	Agreement int8   `json:"agreement" binding:"required,oneof=1"`
	Status    int8   `json:"status" binding:"required,oneof=1 2"` // 1 成功 2 禁用
	Logo      string `json:"logo"`                                // 镜像logo
	Desc      string `json:"desc"`                                // 镜像描述
}

func (MirrorCloudApi) ImageUpdateView(c *gin.Context) {
	cr := middleware.GetBind[ImageUpdateRequest](c)
	var model models.ImageModel
	if err := global.DB.Take(&model, cr.ID).Error; err != nil {
		res.FailWithError(err, c)
		return
	}

	// title不能和现在的重名
	var newModel models.ImageModel
	if err := global.DB.Take(&newModel, "id <> ? and title = ?", cr.ID, cr.Title).Error; err != nil {
		res.FailWithMsg("修改的镜像名称不能重复", c)
		return
	}

	updateData := models.ImageModel{
		Title:     cr.Title,
		Port:      cr.Port,
		Agreement: cr.Agreement,
		Status:    cr.Status,
		Logo:      cr.Logo,
		Desc:      cr.Desc,
	}
	if err := global.DB.Model(&model).Updates(updateData).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("镜像修改成功", c)
}
