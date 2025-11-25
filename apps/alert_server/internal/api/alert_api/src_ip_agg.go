package alert_api

import (
	"alert_server/internal/core"
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

type SrcIpAggRequest struct {
	models.PageInfo
	SrcIp string `form:"srcIp"`
}

type SrcIpAggResponse struct {
	ScrIp         string   `json:"srcIp"`
	Addr          string   `json:"addr"`
	SignatureList []string `json:"signatureList"`
	AttackCount   int      `json:"attackCount"`
	HoneyIpCount  int      `json:"honeyIpCount"`
	NewAttackDate string   `json:"newAttackDate"`
	IsWhite       bool     `json:"isWhite"` // 是否在白名单中
}

type AggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key       string `json:"key"`
		DocCount  int    `json:"doc_count"`
		Signature struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"signature"`
		IpCount struct {
			Value int `json:"value"`
		} `json:"ipCount"`
		MaxDate struct {
			Value         float64 `json:"value"`
			ValueAsString string  `json:"value_as_string"`
		} `json:"maxDate"`
	} `json:"buckets"`
}

type AllAggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}

func (AlertApi) SrcIpAggView(c *gin.Context) {
	cr := middleware.GetBind[SrcIpAggRequest](c)

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

	agg := elastic.NewTermsAggregation()
	agg.Field("srcIp").
		SubAggregation("signature", elastic.NewTermsAggregation().Field("signature.keyword")).
		SubAggregation("ipCount", elastic.NewCardinalityAggregation().Field("destIp")).
		SubAggregation("maxDate", elastic.NewMaxAggregation().Field("timestamp")).
		SubAggregation("page", elastic.NewBucketSortAggregation().
			Sort("maxDate", false).
			From(offset).
			Size(cr.Limit))

	query := elastic.NewBoolQuery()
	// 添加源ip查询条件
	if cr.SrcIp != "" {
		query = query.Filter(elastic.NewTermQuery("srcIp", cr.SrcIp))
	}
	response, err := global.ES.Search(es_model.AlertModel{}.Index()).
		Query(query).
		Aggregation("agg", agg).
		Aggregation("allAgg", elastic.NewTermsAggregation().Field("SrcIp")).
		Size(cr.Limit).
		From(offset).Do(context.Background())
	if err != nil {
		logrus.Errorf("查询告警失败 %s", err)
		res.FailWithMsg("查询告警失败", c)
		return
	}

	var allAggType AllAggType
	err = json.Unmarshal(response.Aggregations["allAgg"], &allAggType)
	if err != nil {
		logrus.Errorf("json解析失败 %s %s", err, response.Aggregations["allAgg"])
		res.FailWithMsg("json解析失败", c)
		return
	}
	count := len(allAggType.Buckets)

	var aggType AggType
	err = json.Unmarshal(response.Aggregations["agg"], &aggType)
	if err != nil {
		logrus.Errorf("json解析失败 %s %s", err, response.Aggregations["agg"])
		res.FailWithMsg("json解析失败", c)
		return
	}

	var list = make([]SrcIpAggResponse, 0)

	var srcIPList []string
	for _, bucket := range aggType.Buckets {
		srcIPList = append(srcIPList, bucket.Key)
	}

	var whiteIPList []models.WhiteIPModel
	global.DB.Find(&whiteIPList, "ip in ?", srcIPList)
	var maps = map[string]bool{}
	for _, model := range whiteIPList {
		maps[model.IP] = true
	}

	for _, bucket := range aggType.Buckets {
		var signatureList []string
		for _, s := range bucket.Signature.Buckets {
			signatureList = append(signatureList, s.Key)
		}
		list = append(list, SrcIpAggResponse{
			ScrIp:         bucket.Key,
			Addr:          core.GetIpAddr(cr.SrcIp),
			SignatureList: signatureList,
			AttackCount:   bucket.DocCount,
			HoneyIpCount:  bucket.IpCount.Value,
			NewAttackDate: bucket.MaxDate.ValueAsString,
			IsWhite:       maps[bucket.Key],
		})
	}

	res.OkWithList(list, int64(count), c)
}
