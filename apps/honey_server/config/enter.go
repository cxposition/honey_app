package config

import "fmt"

type Config struct {
	DB     DB     `yaml:"db"`
	Logger Logger `yaml:"logger"`
	Redis  Redis  `yaml:"redis"`
}

type DB struct {
	DbName          string `yaml:"db_name"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifeTime int    `yaml:"connMaxLifeTime"`
}

func (db *DB) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DbName,
	)
}

type Logger struct {
	Format  string `yaml:"format"`
	Level   string `yaml:"level"`
	AppName string `yaml:"appName"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
