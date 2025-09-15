package vs_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/docker_service"
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

	// 判断有没有这个镜像有没有跑过这个服务
	var service models.ServiceModel
	err = global.DB.Take(&service, "image_id = ?", cr.ImageID).Error
	if err != nil {
		res.FailWithMsg("镜像已运行虚拟服务", c)
		return
	}
	// 使用docker命令运行容器
	// docker run -d --name xxx -p
	ip := "10.2.0.10"
	networkName := "honey-hy"
	containerName := "hy_" + image.ImageName
	containerID, err := docker_service.RunContainer(containerName, networkName, ip, fmt.Sprintf("%s:%s", image.ImageName, image.Tag))
	if err != nil {
		logrus.Errorf("保存虚拟服务记录失败 %s", err)
		res.FailWithMsg("保存虚拟服务记录失败", c)
		return
	}
	command := fmt.Sprintf("docker run -d --network honey-hy --ip %s --name %s %s:%s",
		ip, image.ImageName, image.ImageName, image.Tag)
	fmt.Println(command)
	var model = models.ServiceModel{
		Title:       image.Title,
		Agreement:   image.Agreement,
		ImageID:     image.ID,
		IP:          ip,
		Port:        image.Port,
		Status:      1,
		ContainerID: containerID,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		logrus.Errorf("保存虚拟服务记录失败 %s", err)
		res.FailWithMsg("保存虚拟服务记录失败", c)
	}
	res.OkWithMsg("创建虚拟服务成功", c)
	return
}
