package flags

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/models"
)

func Migrate() {
	err := global.DB.AutoMigrate(
		&models.PortModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
