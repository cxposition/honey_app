package log_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"
)

type LogListRequest struct {
	models.PageInfo
	Type int8   `form:"type"` // 1 表示登陆日志
	IP   string `form:"ip"`
	Addr string `form:"addr"`
}

func (l *LogApi) LogListView(c *gin.Context) {
	cr := middleware.GetBind[LogListRequest](c)
	list, count, _ := common_service.QueryList[models.LogModel](models.LogModel{
		Type: cr.Type,
		IP:   cr.IP,
		Addr: cr.Addr,
	}, common_service.Request{
		Likes:    []string{"username"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
