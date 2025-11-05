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

type NetUseIPListRequest struct {
	models.IDRequest
	models.PageInfo
}

func (NetApi) NetUseIPListView(c *gin.Context) {
	cr := middleware.GetBind[NetUseIPListRequest](c)

	var model models.NetModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 若无可用范围则自动生成
	if model.CanUseHoneyIPRange == "" {
		_, ipNet, err := net.ParseCIDR(model.Subnet())
		if err != nil {
			res.FailWithMsg("无效的子网格式", c)
			return
		}
		logrus.Infof("ipNet:%+v, ipNet.IP:%+v", ipNet, ipNet.IP)
		startIP := ip.IncrementIP(ipNet.IP)
		endIP := ip.DecrementIP(ip.BroadcastIP(ipNet))
		model.CanUseHoneyIPRange = ip.FormatIPRange(startIP, endIP)
	}

	ipList, err := ip.ParseIPRange(model.CanUseHoneyIPRange)
	if err != nil {
		res.FailWithMsg("IP 范围解析失败", c)
		return
	}

	// 查询已使用 IP
	var filterIPList1, filterIPList2 []string
	global.DB.Model(models.HostModel{}).Where("net_id = ?", cr.ID).Select("ip").Scan(&filterIPList1)
	global.DB.Model(models.HoneyIpModel{}).Where("net_id = ?", cr.ID).Select("ip").Scan(&filterIPList2)

	// 合并已使用 IP
	usedIPs := make(map[string]struct{})
	for _, _ip := range append(filterIPList1, filterIPList2...) {
		usedIPs[_ip] = struct{}{}
	}

	// 生成可用IP列表（过滤保留IP与已用IP）
	var availableIPs []string
	for _, addr := range ipList {
		if strings.HasSuffix(addr, ".1") || strings.HasSuffix(addr, ".2") || strings.HasSuffix(addr, ".255") {
			continue
		}
		if _, exists := usedIPs[addr]; !exists {
			availableIPs = append(availableIPs, addr)
		}
	}

	// ---- 前缀匹配过滤 ----
	if cr.Key != "" {
		key := strings.TrimSpace(cr.Key)
		var filtered []string
		for _, addr := range availableIPs {
			if strings.HasPrefix(addr, key) {
				filtered = append(filtered, addr)
			}
		}
		availableIPs = filtered
	}

	// ---- 分页逻辑 ----
	page := cr.Page
	limit := cr.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 254 {
		limit = 20
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(availableIPs) {
		start = len(availableIPs)
	}
	if end > len(availableIPs) {
		end = len(availableIPs)
	}
	pagedIPs := availableIPs[start:end]

	// 返回结果
	res.OkWithData(NetUseIPListResponse{
		Total:              len(ipList),
		Used:               len(ipList) - len(availableIPs),
		UseIPList:          pagedIPs,
		CanUseHoneyIPRange: model.CanUseHoneyIPRange,
	}, c)
}
