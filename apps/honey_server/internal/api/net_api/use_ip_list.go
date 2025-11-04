package net_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"net"
	"strings"
)

type NetUseIPListResponse struct {
	Total              int      `json:"total"`
	Used               int      `json:"used"`
	UseIPList          []string `json:"useIPList"`
	CanUseHoneyIPRange string   `json:"canUseHoneyIPRange"` // 能够使用的诱捕ip范围
}

func (NetApi) NetUseIPListView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NetModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	if model.CanUseHoneyIPRange == "" {
		// 算起始ip
		_, ipNet, err := net.ParseCIDR(model.Subnet())
		if err != nil {
			res.FailWithMsg("无效的子网格式", c)
			return
		}

		logrus.Infof("ipNet:%+v, ipNet.IP:%+v", ipNet, ipNet.IP)

		// 计算可用IP范围（排除网络地址和广播地址）
		startIP := ip.IncrementIP(ipNet.IP)
		endIP := ip.DecrementIP(ip.BroadcastIP(ipNet))
		model.CanUseHoneyIPRange = ip.FormatIPRange(startIP, endIP)
	}
	ipList, err := ip.ParseIPRange(model.CanUseHoneyIPRange)
	// 排查ip
	var filterIPList1, filterIPList2 []string
	global.DB.Model(models.HostModel{}).Where("net_id = ?", cr.ID).Select("ip").Scan(&filterIPList1)
	global.DB.Model(models.HoneyIpModel{}).Where("net_id = ?", cr.ID).Select("ip").Scan(&filterIPList2)

	// 合并已使用IP列表
	usedIPs := make(map[string]struct{})
	for _, _ip := range filterIPList1 {
		usedIPs[_ip] = struct{}{}
	}
	for _, _ip := range filterIPList2 {
		usedIPs[_ip] = struct{}{}
	}

	var availableIPs []string
	for _, addr := range ipList {
		// 跳过网络中常见的保留IP
		if strings.HasSuffix(addr, ".1") || strings.HasSuffix(addr, ".2") || strings.HasSuffix(addr, ".255") {
			continue
		}
		if _, exists := usedIPs[addr]; !exists {
			availableIPs = append(availableIPs, addr)
		}
	}

	// 返回结果
	res.OkWithData(NetUseIPListResponse{
		Total:              len(ipList),
		Used:               len(ipList) - len(availableIPs),
		UseIPList:          availableIPs,
		CanUseHoneyIPRange: model.CanUseHoneyIPRange,
	}, c)
}
