package flags

import (
	"alert_server/internal/global"
	"github.com/sirupsen/logrus"
)

func Migrate() {
	err := global.DB.AutoMigrate()
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
