package vs_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"
)

type VsListRequest struct {
	models.PageInfo
	Port  int    `form:"port"`
	IP    string `form:"IP"`
	Title string `form:"title"`
}

func (VsApi) VsListView(c *gin.Context) {
	cr := middleware.GetBind[VsListRequest](c)
	list, count, _ := common_service.QueryList(models.ServiceModel{
		Port:  cr.Port,
		IP:    cr.IP,
		Title: cr.Title,
	}, common_service.RequestList{
		Likes:    []string{"title"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
