package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image_server/internal/utils/cmd"
	"os"
	"path"
)

type ImageModel struct {
	Model
	ImageName     string         `gorm:"size:64" json:"imageName"`
	Title         string         `gorm:"size:64" json:"title"`
	Port          int            `json:"port"`
	DockerImageID string         `gorm:"size:32" json:"dockerImageID"`
	ServiceList   []ServiceModel `gorm:"foreignKey:ImageID" json:"-"` // 关联的虚拟服务列表
	Tag           string         `gorm:"size:32" json:"tag"`
	Agreement     int8           `json:"agreement"`                 // 协议
	ImagePath     string         `gorm:"size:256" json:"imagePath"` // 镜像文件
	Status        int8           `json:"status"`                    //
	Logo          string         `gorm:"size:256" json:"logo"`      // 镜像logo
	Desc          string         `gorm:"size:512" json:"desc"`      // 镜像描述
}

func (i *ImageModel) BeforeDelete(tx *gorm.DB) error {
	// 删除docker镜像
	command := fmt.Sprintf("docker rmi %s", i.DockerImageID)

	if err := cmd.Cmd(command); err != nil {
		logrus.Errorf("删除镜像失败 %s", err)
		return err
	}

	// 删除镜像文件
	logrus.Infof("删除镜像文件 %s", i.ImagePath)
	currentPath, _ := os.Getwd()
	if err := os.Remove(path.Join(currentPath, i.ImagePath)); err != nil {
		return err
	}
	return nil
}
