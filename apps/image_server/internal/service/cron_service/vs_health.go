package cron_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/service/docker_service"
)

func VsHealth() {
	logrus.Infof("获取前缀 %s 的容器状态", global.Config.VsNet.Prefix)
	allContainers, err := docker_service.PrefixContainerStatus(global.Config.VsNet.Prefix)
	if err != nil {
		logrus.Errorf("容器状态检测失败 %s", err)
		return
	}

	var list []models.ServiceModel
	global.DB.Find(&list)
	var containerMap = map[string]*models.ServiceModel{}
	for _, model := range list {
		containerMap[model.ContainerID] = &model
	}

	for _, container := range allContainers {
		containerID := container.ID[:12]
		logrus.Infof("containerID:%s", containerID)
		model, ok := containerMap[containerID]
		logrus.Infof("ok:%+v", ok)
		if !ok {
			continue
		}
		var newModel models.ServiceModel
		var isUpdate bool
		if container.State == "running" && model.Status != 1 {
			// 我们这边是不正常的，但是实际是正常的
			newModel.Status = 1
			newModel.ErrorMsg = ""
			isUpdate = true
		}
		if container.State != "running" && model.Status == 1 {
			// 我们这边是正常的，但是实际是不正常的
			newModel.Status = 2
			newModel.ErrorMsg = fmt.Sprintf("%s(%s)", container.State, container.Status)
			isUpdate = true
		}

		logrus.Infof("isUpdate:%+v", isUpdate)

		if isUpdate {
			logrus.Infof("%s 容器存在状态修改 %s => %s", model.ContainerName, model.State(), container.State)
			global.DB.Model(model).Updates(map[string]any{
				"status":    newModel.Status,
				"error_msg": newModel.ErrorMsg,
			})
		}
	}
}
