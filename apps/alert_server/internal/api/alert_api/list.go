package alert_api

import (
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"alert_server/internal/middleware"
	"alert_server/internal/models"
	"alert_server/internal/utils/res"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ListRequest struct {
	models.PageInfo
	SrcIp       string `form:"srcIp"`       // 源ip
	DestIp      string `form:"destIp"`      // 目标ip
	DestPort    int    `form:"destPort"`    // 目标端口
	ServiceName string `form:"serviceName"` // 服务名称
	Signature   string `form:"signature"`   // 攻击类型
	Level       int    `form:"level"`       // 告警级别
	StartTime   string `form:"startTime"`   // 开始时间 (2025-01-01 00:00:00)
	EndTime     string `form:"endTime"`     // 结束时间 (2025-01-02 00:00:00)
}

func (AlertApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)

	// 分页处理
	if cr.Limit <= 0 {
		cr.Limit = 10
	}
	if cr.Limit > 20 {
		cr.Limit = 10
	}
	if cr.Page <= 0 {
		cr.Page = 1
	}

	offset := (cr.Page - 1) * cr.Limit

	query := elastic.NewBoolQuery()

	if cr.SrcIp != "" {
		query.Must(elastic.NewTermQuery("srcIp", cr.SrcIp))
	}

	if cr.DestIp != "" {
		query.Must(elastic.NewTermQuery("destIP", cr.DestIp))
	}

	if cr.DestPort != 0 {
		query.Must(elastic.NewTermQuery("destPort", cr.DestPort))
	}

	if cr.ServiceName != "" {
		query.Must(elastic.NewTermQuery("serviceName.keyword", cr.ServiceName))
	}

	if cr.Signature != "" {
		query.Must(elastic.NewTermQuery("signature.keyword", cr.Signature))
	}

	if cr.Level != 0 {
		query.Must(elastic.NewTermQuery("level", cr.Level))
	}

	if cr.StartTime != "" && cr.EndTime != "" {
		query.Must(elastic.NewRangeQuery("timestamp").
			Gte(cr.StartTime).
			Lte(cr.EndTime))
	}

	// 查询 ES
	response, err := global.ES.Search(es_model.AlertModel{}.Index()).
		Query(query).
		Size(cr.Limit).
		From(offset).
		Sort("timestamp", false).
		Do(context.Background())

	if err != nil {
		logrus.Errorf("告警查询失败 %s", err)
		res.FailWithMsg("告警查询失败", c)
		return
	}

	count := response.Hits.TotalHits.Value

	// 解析数据
	var list []es_model.AlertModel
	for _, hit := range response.Hits.Hits {
		var data es_model.AlertModel
		err = json.Unmarshal(hit.Source, &data)
		if err != nil {
			logrus.Errorf("json解析失败 %s %s %s", err, hit.Source, hit.Id)
			continue
		}
		data.ID = hit.Id
		list = append(list, data)
	}

	res.OkWithList(list, count, c)
}
