package es_model

import (
	"alert_server/internal/global"
	_ "embed"
)

type AlertModel struct {
	NodeUid          string `json:"nodeUid"`
	SrcIp            string `json:"srcIp"`
	SrcPort          int    `json:"srcPort"`
	DestIP           string `json:"destIP"`
	DestPort         int    `json:"destPort"`
	Timestamp        string `json:"timestamp"` // 年月日，时分秒的时间
	Signature        string `json:"signature"`
	Level            int8   `json:"level"` // 告警级别
	HttpResponseBody string `json:"httpResponseBody"`
	Payload          string `json:"payload"`
}

func (alert AlertModel) Index() string {
	return global.Config.Alert.AlertIndex
}

//go:embed alert_mapping.json
var alertMapping string

func (alert AlertModel) Mappings() string {
	return alertMapping
}
