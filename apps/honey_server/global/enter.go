package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"honey_app/apps/honey_server/config"
)

var (
	DB     *gorm.DB
	Config *config.Config
	Log    *logrus.Entry
)
