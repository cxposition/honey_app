package flags

import (
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"context"
	"github.com/sirupsen/logrus"
)

func ESIndex() {
	index := es_model.AlertModel{}.Index()
	ok, err := global.ES.IndexExists(index).Do(context.Background())
	if err != nil {
		logrus.Errorf("获取索引错误 %s", err)
		return
	}
	if ok {
		logrus.Infof("存在索引 删除索引 %s", index)
		global.ES.DeleteIndex(index).Do(context.Background())
	}
	logrus.Infof("创建索引 %s", index)
	response, err := global.ES.CreateIndex(index).Body(es_model.AlertModel{}.Mappings()).Do(context.Background())
	if err != nil {
		logrus.Errorf("创建索引错误 %s", err)
		return
	}
	logrus.Infof("创建索引成功 %v", response)
}
