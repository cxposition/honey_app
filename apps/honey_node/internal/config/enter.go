package config

import "fmt"

type Config struct {
	Logger            Logger   `yaml:"logger"`
	System            System   `yaml:"system"`
	FilterNetworkList []string `yaml:"filterNetworkList"`
	MQ                MQ       `yaml:"mq"`
	DB                DB       `yaml:"db"`
}

type DB struct {
	DbName          string `yaml:"db_name"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifeTime int    `yaml:"connMaxLifeTime"`
}

//func (db *DB) DSN() string {
//	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//		db.User,
//		db.Password,
//		db.Host,
//		db.Port,
//		db.DbName,
//	)
//}

type Logger struct {
	Format  string `yaml:"format"`
	Level   string `yaml:"level"`
	AppName string `yaml:"appName"`
}

type System struct {
	GrpcManageAddr string `yaml:"grpcManageAddr"`
	Network        string `yaml:"network"`
	Uid            string `yaml:"uid"`
	EvePath        string `yaml:"evePath"`
}

type MQ struct {
	User                 string `yaml:"user"`
	Password             string `yaml:"password"`
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
	CreateIpExchangeName string `yaml:"createIpExchangeName"`
	DeleteIpExchangeName string `yaml:"deleteIpExchangeName"`
	BindPortExchangeName string `yaml:"bindPortExchangeName"`
	Ssl                  bool   `yaml:"ssl"`
	ClientCertificate    string `yaml:"clientCertificate"`
	ClientKey            string `yaml:"clientKey"`
	CaCertificate        string `yaml:"caCertificate"`
	AlertTopic           string `yaml:"alertTopic"`
}

func (m MQ) Addr() string {
	if m.Ssl {
		return fmt.Sprintf("amqps://%s:%s@%s:%d/",
			m.User,
			m.Password,
			m.Host,
			m.Port,
		)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		m.User,
		m.Password,
		m.Host,
		m.Port,
	)
}
