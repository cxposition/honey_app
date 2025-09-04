package flags

import (
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/internal/global"
	models2 "honey_app/apps/honey_server/internal/models"
)

func Migrate() {
	err := global.DB.AutoMigrate(
		&models2.HoneyIpModel{},
		&models2.HoneyPortModel{},
		&models2.HostModel{},
		&models2.HostTemplateModel{},
		&models2.ImageModel{},
		&models2.LogModel{},
		&models2.MatrixTemplateModel{},
		&models2.NetModel{},
		&models2.NodeModel{},
		&models2.NodeNetworkModel{},
		&models2.ServiceModel{},
		&models2.UserModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
