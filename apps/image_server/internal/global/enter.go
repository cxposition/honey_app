package global

import (
	"github.com/docker/docker/client"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"image_server/internal/config"
)

var (
	Version   = "v1.0.1"
	Commit    = "a5f28b47b9"
	BuildTime = "2025-08-21 21:30:05"
)

var (
	DB           *gorm.DB
	Redis        *redis.Client
	Config       *config.Config
	Log          *logrus.Entry
	DockerClient *client.Client
)
