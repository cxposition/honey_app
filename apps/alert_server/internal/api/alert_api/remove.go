package alert_api

import (
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"alert_server/internal/middleware"
	"alert_server/internal/utils/res"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type RemoveRequest struct {
	ID    string `json:"id"`
	SrcIp string `json:"srcIp"`
}

func (AlertApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[RemoveRequest](c)
	var idList []string
	if cr.ID != "" {
		idList = append(idList, cr.ID)
	}
	if cr.SrcIp != "" {
		response, err := global.ES.Search(es_model.AlertModel{}.Index()).
			Query(elastic.NewTermQuery("srcIp", cr.SrcIp)).
			Size(10000).
			Do(context.Background())
		if err != nil {
			logrus.Errorf("查询失败: %v", err)
			res.FailWithMsg("查询失败", c)
			return
		}

		for _, hit := range response.Hits.Hits {
			idList = append(idList, hit.Id)
		}
	}

	if len(idList) == 0 {
		res.FailWithMsg("不存在的告警记录", c)
		return
	}

	index := es_model.AlertModel{}.Index()

	bulk := global.ES.Bulk()
	for _, id := range idList {
		bulk.Add(elastic.NewBulkDeleteRequest().Index(index).Id(id))
	}

	resp, err := bulk.Refresh("true").Do(context.Background())
	if err != nil {
		logrus.Errorf("ES批量删除失败: %s", err)
		res.FailWithMsg("批量删除失败", c)
		return
	}

	deleted := 0
	failed := 0
	for _, items := range resp.Items {
		for _, item := range items {
			if item.Error == nil && item.Status >= 200 && item.Status < 300 {
				deleted++
			} else {
				failed++
				logrus.Errorf("删除失败 id=%s status=%d err=%v", item.Id, item.Status, item.Error)
			}
		}
	}

	msg := fmt.Sprintf("删除请求 %d 条，成功 %d 条，失败 %d 条", len(idList), deleted, failed)
	res.OkWithMsg(msg, c)
}
