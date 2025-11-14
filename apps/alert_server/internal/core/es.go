package core

import (
	"alert_server/internal/global"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func InitES() *elastic.Client {
	cfg := global.Config.ES
	sniffOpt := elastic.SetSniff(false)
	c, err := elastic.NewClient(
		elastic.SetURL(cfg.Addr),
		sniffOpt,
		elastic.SetBasicAuth(cfg.UserName, cfg.Password),
	)
	if err != nil {
		logrus.Fatalf("es连接失败 %s", err.Error())
	}
	logrus.Infof("es连接成功")
	return c
}
