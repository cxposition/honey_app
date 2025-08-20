package global

import (
	"gorm.io/gorm"
	"honey_app/apps/honey_server/config"
)

var (
	DB *gorm.DB
)

var Config *config.Config
