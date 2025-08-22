package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"honey_app/apps/honey_server/config"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Log    *logrus.Entry
)

var (
	Version   = "v1.0.1"
	Commit    = "a5f28b47b9"
	BuildTime = "2025-08-21 21:30:05"
)
