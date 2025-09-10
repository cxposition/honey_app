package mirror_cloud_api

import (
	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

func (MirrorCloudApi) ImageRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	log := middleware.GetLog(c)
	var model models.ImageModel
	if err := global.DB.Preload("ServiceList").Take(&model, cr.ID).Error; err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}
	if len(model.ServiceList) > 0 {
		res.FailWithMsg("镜像存在虚拟服务，请先删除关联的虚拟服务", c)
		return
	}
	if err := global.DB.Delete(&model).Error; err != nil {
		log.Errorf("删除镜像失败: %v", err)
		res.FailWithMsg("删除镜像失败", c)
		return
	}
	res.OkWithMsg("删除镜像成功", c)
}
