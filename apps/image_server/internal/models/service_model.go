package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image_server/internal/utils/cmd"
)

type ServiceModel struct {
	Model
	Title         string     `json:"title"`
	Agreement     int8       `json:"agreement"`
	ImageID       uint       `json:"imageID"`
	ImageModel    ImageModel `gorm:"foreignKey:ImageID" json:"-"`
	IP            string     `json:"ip"`
	Port          int        `json:"port"`
	Status        int8       `json:"status"`
	ErrorMsg      string     `json:"errorMsg"`
	HoneyIPCount  int        `json:"honeyIPCount"`
	ContainerID   string     `json:"containerID"`                  // 容器ID
	ContainerName string     `gorm:"size:32" json:"containerName"` // 容器名
}

func (s *ServiceModel) State() string {
	switch s.Status {
	case 1:
		return "running"

	}
	return "error"
}

func (s *ServiceModel) BeforeDelete(tx *gorm.DB) error {
	// 判断是否有关联的端口转发
	var count int64
	tx.Model(&HoneyPortModel{}).Where("service_id = ?", s.ID).Count(&count)
	if count > 0 {
		return errors.New("存在端口转发,不能删除虚拟服务")
	}

	command := fmt.Sprintf("docker rm -f %s", s.ContainerName)
	err := cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("删除容器失败 %s", err)
		return err
	}
	return nil

}
