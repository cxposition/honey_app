package mq_service

import (
	"alert_server/internal/core"
	"alert_server/internal/es_model"
	"alert_server/internal/global"
	"alert_server/internal/models"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func RevAlertMQ(ch *amqp.Channel) {
	cfg := global.Config.Alert
	msgs, err := ch.Consume(cfg.AlertTopic, "", true, false, false, false, nil)
	if err != nil {
		logrus.Fatalf("消费mq队列失败 %s", err)
	}
	if err != nil {
		logrus.Fatalf("消费mq队列失败 %s", err)
	}

	for d := range msgs {
		var data es_model.AlertModel
		err = json.Unmarshal(d.Body, &data)
		if err != nil {
			logrus.Errorf("消息格式解析失败 %s %s", err, d.Body)
			continue
		}

		logrus.Infof("%s %s => %s:%d", data.Signature, data.SrcIp, data.DestIP, data.DestPort)
		// 查ip是否在白名单中
		var whiteModel models.WhiteIPModel
		global.DB.Find(&whiteModel, "ip = ?", data.SrcIp)
		if whiteModel.ID != 0 {
			logrus.Warnf("告警消息 在白名单中")
			continue
		}

		// 查虚拟服务
		var hpModel models.HoneyPortModel
		global.DB.Preload("ServiceModel").Find(&hpModel, "ip = ? and port = ?", data.DestIP, data.DestPort)
		if hpModel.ID != 0 {
			data.ServiceID = hpModel.ServiceID
			data.ServiceName = hpModel.ServiceModel.Title
		}

		addr := core.GetIpAddr(data.SrcIp)
		data.Addr = addr

		response, err := global.ES.Index().Index(data.Index()).BodyJson(data).Do(context.Background())
		if err != nil {
			logrus.Errorf("告警消息入库失败 %s %s", err, d.Body)
			continue
		}
		logrus.Infof("告警消息入库成功 %s", response.Id)

	}
}
