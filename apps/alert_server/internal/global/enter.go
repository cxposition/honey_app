package global

import (
	"alert_server/internal/config"
	"github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Log    *logrus.Entry
	MQConn *amqp.Connection
)

var (
	Version   = "v1.0.1"
	Commit    = "a5f28b47b9"
	BuildTime = "2025-08-21 21:30:05"
	ES        *elastic.Client
)
