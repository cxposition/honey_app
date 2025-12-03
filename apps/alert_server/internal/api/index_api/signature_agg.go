package index_api

import (
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"alert_server/internal/utils/res"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type SignatureAggResponse struct {
	Signature string `json:"signature"`
	Count     int    `json:"count"`
}

func (IndexApi) SignatureAggView(c *gin.Context) {

	agg := elastic.NewTermsAggregation().
		Field("signature.keyword")

	resp, err := global.ES.Search(es_model.AlertModel{}.Index()).
		Size(0).
		Aggregation("signatureAgg", agg).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("Signature聚合查询失败: %v", err)
		res.FailWithMsg("查询失败", c)
		return
	}

	terms, found := resp.Aggregations.Terms("signatureAgg")
	if !found || terms == nil {
		res.OkWithData([]SignatureAggResponse{}, c)
		return
	}

	var list = make([]SignatureAggResponse, 0, len(terms.Buckets))
	for _, b := range terms.Buckets {
		list = append(list, SignatureAggResponse{
			Signature: fmt.Sprintf("%v", b.Key),
			Count:     int(b.DocCount),
		})
	}

	res.OkWithData(list, c)
}
