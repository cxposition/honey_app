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
	SrcIp       string `form:"srcIp"`
	DestIp      string `form:"destIp"`
	DestPort    int    `form:"destPort"`
	ServiceName string `form:"serviceName"`
	Signature   string `form:"signature"`
	Level       int    `form:"level"`
	StartTime   string `form:"startTime"`
	EndTime     string `form:"endTime"`
}

func (AlertApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[ListRequest](c)
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
	response, err := global.ES.Search(es_model.AlertModel{}.Index()).
		Query(query).
		Size(cr.Limit).
		From(offset).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("告警查询失败 %s", err)
		res.FailWithMsg("告警查询失败 %s", c)
		return
	}
	count := response.Hits.TotalHits.Value

	var list = make([]es_model.AlertModel, 0)
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
