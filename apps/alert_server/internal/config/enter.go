package config

import "fmt"

type Config struct {
	DB        DB       `yaml:"db"`
	Logger    Logger   `yaml:"logger"`
	Redis     Redis    `yaml:"redis"`
	System    System   `yaml:"system"`
	Jwt       Jwt      `yaml:"jwt"`
	WhiteList []string `yaml:"whiteList"`
	MQ        MQ       `yaml:"mq"`
	ES        ES       `yaml:"es"`
	Alert     Alert    `yaml:"alert"`
}

type ES struct {
	Addr     string `yaml:"addr"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
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

type System struct {
	WebAddr string `yaml:"webAddr"`
}
type Jwt struct {
	Expires int    `yaml:"expires"` // 单位为秒
	Issuer  string `yaml:"issuer"`
	Secret  string `yaml:"secret"`
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

type Alert struct {
	AlertIndex string `yaml:"alertIndex"`
}
