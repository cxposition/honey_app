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

type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIPID"`
	IP        string `json:"ip"`
	Mask      int8   `json:"mask"`
	Network   string `json:"network"` // 基于哪个接口创建
	LogID     string `json:"logID"`   // 日志id
}

func CreateIpExchange(msg string) error {
	var req CreateIPRequest
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		logrus.Errorf("json解析失败 %s %s", err, msg)
		return nil // 不重发消息
	}

	/*
		ip link add mc_12 link ens33 type macvlan mode bridge
		ip link set mc_12 up
		ip addr add 192.168.177.130/24 dev mc_12
	*/

	var errMsg string
	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID)
	err = cmd.Cmd(fmt.Sprintf("ip link add %s link %s type macvlan mode bridge", linkName, req.Network))
	if err != nil {
		errMsg = err.Error()
	}
	err = cmd.Cmd(fmt.Sprintf("ip link set %s up", linkName))
	if err != nil {
		errMsg = err.Error()
	}
	err = cmd.Cmd(fmt.Sprintf("ip addr add %s/%d dev %s", req.IP, req.Mask, linkName))
	if err != nil {
		errMsg = err.Error()
	}

	// 获取mac地址
	mac, err := cmd.CommandWithOut("ip link show ens33 | awk '/ether/ {print $2}' | tr -d '\\n'\n")
	if err != nil {
		errMsg = err.Error()
	}
	// 调用grpc方法上报状态

	response, err := global.GrpcClient.StatusCreateIP(context.Background(), &node_rpc.StatusCreateRequest{
		HoneyIPID: uint32(req.HoneyIPID),
		ErrMsg:    errMsg,
		Network:   linkName,
		Mac:       mac,
	})
	if err != nil {
		logrus.Errorf("调用grpc方法上报状态失败 %s", err)
		return err // 重新发送消息
	}

	logrus.Infof("上报管理状态成功:%v", response)
	return nil
}
