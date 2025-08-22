package core

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"honey_app/apps/honey_server/global"
	"sync"
	"time"
)

var db *gorm.DB
var dbOnce sync.Once

func InitDB() (database *gorm.DB) {
	cfg := global.Config.DB
	dsn := cfg.DSN()
	dialector := mysql.Open(dsn)
	database, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 不生成实体外键
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}
	// 配置连接池
	sqlDB, err := database.DB()
	if err != nil {
		logrus.Fatalf("获取数据库连接失败 %s", err)
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}
	// 设置连接池
	var maxIdleConns, maxOpenConns, connMaxLifetime int
	if cfg.MaxIdleConns == 0 {
		maxIdleConns = 10
	} else {
		maxIdleConns = cfg.MaxIdleConns
	}
	if cfg.MaxOpenConns == 0 {
		maxOpenConns = 100
	} else {
		maxOpenConns = cfg.MaxOpenConns
	}
	if cfg.ConnMaxLifeTime == 0 {
		connMaxLifetime = 10000
	} else {
		connMaxLifetime = cfg.ConnMaxLifeTime
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	logrus.Infof("mysql连接成功")
	return
}

func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		db = InitDB()
	})
	return db
}
