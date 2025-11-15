package main

import (
	"alert_server/internal/core"
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	global.Config = core.ReadConfig()
	core.SetLogDefault()
	global.ES = core.InitES()
	remove()
}

func create() {
	model := es_model.AlertModel{
		NodeUid:          "xxx02",
		SrcIp:            "192.168.20.5",
		SrcPort:          3303,
		DestIP:           "192.168.100.5",
		DestPort:         200,
		Timestamp:        time.Now().Format(time.DateTime),
		Signature:        "curl请求",
		Level:            2,
		HttpResponseBody: "xxx",
		Payload:          "xxxx",
		ServiceID:        1,
		ServiceName:      "web服务",
	}

	response, err := global.ES.Index().Index(model.Index()).BodyJson(model).Do(context.Background())
	if err != nil {
		logrus.Errorf("es数据写入失败 %s", err)
		return
	}
	logrus.Infof("数据写入成功 %#v", response)
}

func list() {
	limit := 1
	page := 1
	offset := (page - 1) * limit
	query := elastic.NewBoolQuery()
	res, err := global.ES.Search(es_model.AlertModel{}.Index()).Query(query).Size(limit).From(offset).Do(context.Background())
	if err != nil {
		logrus.Errorf("es查询失败 %s", err)
		return
	}
	fmt.Println(res.Hits.TotalHits.Value)

	for _, hit := range res.Hits.Hits {
		fmt.Println(hit.Id, string(hit.Source))
	}
}

func remove() {
	deleteResponse, err := global.ES.
		Delete().                             // 调用删除接口
		Index(es_model.AlertModel{}.Index()). // 指定索引
		Id("SutShZoB5MOH293S1wOZ").           // 指定id
		Refresh("true").                      // 是否立即生效
		Do(context.Background())              // 执行
	if err != nil {
		logrus.Errorf("es数据删除失败 %s", err)
		return
	}
	logrus.Infof("数据删除成功 %#v", deleteResponse)
}
