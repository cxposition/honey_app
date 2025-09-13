package vs_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/cmd"
	"image_server/internal/utils/res"
)

type VsCreateRequest struct {
	ImageID int `json:"imageID" binding:"required"`
}

func (VsApi) VsCreateView(c *gin.Context) {
	cr := middleware.GetBind[VsCreateRequest](c)
	var image models.ImageModel
	err := global.DB.Take(&image, cr.ImageID).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if image.Status == 2 {
		res.FailWithMsg("镜像已停用，请重新选择镜像", c)
		return
	}

	// 使用docker命令运行容器
	// docker run -d --name xxx -p
	ip := "10.2.0.10"
	command := fmt.Sprintf("docker run -d --network honey-hy --ip %s --name %s %s:%s",
		ip, image.ImageName, image.ImageName, image.Tag)
	err = cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}
	var model = models.ServiceModel{
		Title:       image.Title,
		Agreement:   image.Agreement,
		ImageID:     image.ID,
		IP:          ip,
		Port:        image.Port,
		Status:      1,
		ContainerID: "",
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		logrus.Errorf("保存虚拟服务记录失败 %s", err)
		res.FailWithMsg("保存虚拟服务记录失败", c)
	}
	res.OkWithMsg("创建虚拟服务成功", c)
	return
}
