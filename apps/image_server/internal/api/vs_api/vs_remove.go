package vs_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

func (VsApi) VsRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c)
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", cr.IDList)
	if len(serviceList) == 0 {
		res.FailWithMsg("虚拟服务不存在", c)
		return
	}
	result := global.DB.Delete(&serviceList)
	successCount := result.RowsAffected
	err := result.Error
	if err != nil {
		res.FailWithMsg("删除虚拟服务失败", c)
		return
	}
	msg := fmt.Sprintf("删除虚拟服务成功, 删除数量: %d", successCount)
	res.OkWithMsg(msg, c)
}
