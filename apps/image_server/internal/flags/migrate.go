package flags

import (
	"github.com/sirupsen/logrus"
	"image_server/internal/global"
	"image_server/internal/models"
)

func Migrate() {
	err := global.DB.AutoMigrate(
		&models.HostTemplateModel{},
		&models.ImageModel{},
		&models.MatrixTemplateModel{},
		&models.ServiceModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
