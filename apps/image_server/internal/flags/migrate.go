package flags

import (
	"github.com/sirupsen/logrus"
	"honey_app/apps/image_server/internal/global"
	models2 "honey_app/apps/image_server/internal/models"
)

func Migrate() {
	err := global.DB.AutoMigrate(
		&models2.HostTemplateModel{},
		&models2.ImageModel{},
		&models2.MatrixTemplateModel{},
		&models2.ServiceModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
