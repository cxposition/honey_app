package white_ip_api

import (
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"alert_server/internal/service/common_service"
	"alert_server/internal/utils/res"
	"github.com/gin-gonic/gin"
)

func (WhiteIPApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[models.PageInfo](c)
	list, count, _ := common_service.QueryList(models.WhiteIPModel{}, common_service.ListRequest{
		Likes:    []string{"ip", "notice"},
		Sort:     "created_at desc",
		PageInfo: cr,
	})
	res.OkWithList(list, count, c)
}
