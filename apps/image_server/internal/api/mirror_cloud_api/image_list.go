package mirror_cloud_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"
)

type ImageListRequest struct {
	models.PageInfo
}

func (MirrorCloudApi) ImageListView(c *gin.Context) {
	cr := middleware.GetBind[ImageListRequest](c)
	list, count, _ := common_service.QueryList(models.ImageModel{},
		common_service.RequestList{
			Debug:    true,
			Likes:    []string{"title", "image_name"}, // username like req.Key
			PageInfo: cr.PageInfo,
			Sort:     "created_at desc",
		})
	res.OkWithList(list, count, c)
}
