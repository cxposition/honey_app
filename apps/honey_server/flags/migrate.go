package flags

import (
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/models"
)

func Migrate() {
	err := global.DB.AutoMigrate(
		&models.HoneyIpModel{},
		&models.HoneyPortModel{},
		&models.HostModel{},
		&models.HostTemplateModel{},
		&models.ImageModel{},
		&models.LogModel{},
		&models.MatrixTemplateModel{},
		&models.NetModel{},
		&models.NodeModel{},
		&models.NodeNetworkModel{},
		&models.ServiceModel{},
		&models.UserModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
