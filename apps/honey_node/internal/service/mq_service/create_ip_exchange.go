package mq_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/j-keck/arping"
	"github.com/sirupsen/logrus"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/models"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/ip_service"
	"net"
)

type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIPID"`
	IP        string `json:"ip"`
	Mask      int8   `json:"mask"`
	Network   string `json:"network"` // 基于哪个接口创建
	LogID     string `json:"logID"`   // 日志id
	IsTan     bool   `json:"isTan"`
}

func CreateIpExchange(msg string) error {
	log := core.GetLogger()
	var req CreateIPRequest
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil
	}

	// 判断是否是探针ip
	if req.IsTan {
		mac, _ := ip_service.GetMACAddress(req.Network)
		return reportStatus(req.HoneyIPID, req.Network, mac, "", req.LogID)
	}

	// 记录处理开始
	log.WithFields(logrus.Fields{
		"req_data": req,
	}).Info("开始处理创建IP请求")

	// 要去判断这个ip有没有使用
	_mac, _, err := arping.PingOverIfaceByName(net.ParseIP(req.IP), req.Network)
	if err == nil {
		err = fmt.Errorf("创建诱捕ip失败 ip已存在 ip %s mac %s", req.IP, _mac.String())
		log.Error(err)
		return reportStatus(req.HoneyIPID, "", _mac.String(), err.Error(), req.LogID)
	}

	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID)

	mac, err := ip_service.SetIp(ip_service.SetIpRequest{
		Ip:       req.IP,
		Mask:     req.Mask,
		LinkName: linkName,
		Network:  req.Network,
	})
	if err != nil {
		return reportStatus(req.HoneyIPID, linkName, mac, err.Error(), req.LogID)
	}

	// 消息持久化
	global.DB.Create(&models.IpModel{
		Ip:       req.IP,
		Mask:     req.Mask,
		LinkName: linkName,
		Network:  req.Network,
	})

	// 所有步骤成功，上报状态
	return reportStatus(req.HoneyIPID, linkName, mac, "", req.LogID)
}

// 上报状态到管理服务
func reportStatus(honeyIPID uint, network, mac, errMsg string, logID string) error {
	log := core.GetLogger().WithField("logID", logID)
	data := &node_rpc.StatusCreateIPRequest{
		HoneyIPID: uint32(honeyIPID),
		ErrMsg:    errMsg,
		Network:   network,
		Mac:       mac,
		LogID:     logID,
	}
	_, err := global.GrpcClient.StatusCreateIP(context.Background(), data)
	if err != nil {
		log.WithField("error", err).Errorf("上报管理状态失败")
		return err
	}

	log.WithFields(logrus.Fields{
		"data": data,
	}).Infof("上报管理状态成功")

	return nil
}
