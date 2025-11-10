package mq_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/cmd"
)

type DeleteIPRequest struct {
	IpList []IpInfo `json:"ipList"`
	LogID  string   `json:"logID"` // 日志id
}

type IpInfo struct {
	HoneyIPID uint   `json:"honeyIPID"`
	IP        string `json:"ip"`
	Network   string `json:"network"` // 基于哪个接口创建
}

func DeleteIpExchange(msg string) error {
	var req DeleteIPRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("json 解析失败: %v,消息: %s", err, msg)
		return nil
	}
	global.Log.WithFields(logrus.Fields{"req": req}).Infof("删除诱捕ip")

	var idList []uint32
	for _, info := range req.IpList {
		cmd.Cmd(fmt.Sprintf("ip link del %s", info.Network))
		idList = append(idList, uint32(info.HoneyIPID))
	}

	reportDeleteIPStatus(idList)
	return nil
}

// 上报状态到管理服务
func reportDeleteIPStatus(honeyIPIDList []uint32) error {
	response, err := global.GrpcClient.StatusDeleteIP(context.Background(), &node_rpc.StatusDeleteIPRequest{
		HoneyIPIDList: honeyIPIDList,
	})
	if err != nil {
		logrus.Errorf("上报管理状态失败: %v", err)
		return err
	}

	global.Log.WithFields(logrus.Fields{"honeyIPIDList": honeyIPIDList}).Infof("上报管理状态成功: %v", response)
	return nil
}
