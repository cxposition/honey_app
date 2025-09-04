package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"sync"
)

var rdb *redis.Client
var rdbOnce sync.Once

func InitRedis() (client *redis.Client) {
	conf := global.Config.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("连接redis失败 %s", err)
		return
	}
	logrus.Infof("成功连接redis")
	return rdb
}

// GetRedis 这有个坑,这个once被mysql用掉了，所以要新建一个once
func GetRedis() *redis.Client {
	rdbOnce.Do(func() {
		rdb = InitRedis()
	})
	return rdb
}
