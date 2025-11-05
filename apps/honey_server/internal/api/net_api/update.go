package net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"net"
)

type UpdateRequest struct {
	ID                 uint   `json:"id" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Gateway            string `json:"gateway"`
	CanUseHoneyIPRange string `json:"canUseHoneyIpRange"`
}

func (NetApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)
	var model models.NetModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	if cr.Title != model.Title {
		var newNet models.NetModel
		err = global.DB.Take(&newNet, "title = ? and id <> ?", cr.Title, cr.ID).Error
		if err == nil {
			res.FailWithMsg("修改的网络名称不能重复", c)
			return
		}
	}

	if cr.Gateway != "" {
		gateway := net.ParseIP(cr.Gateway)
		if gateway == nil || gateway.To4() == nil {
			res.FailWithMsg("网关ip格式错误", c)
			return
		}
		to4 := gateway.To4()
		if to4 == nil {
			res.FailWithMsg("网关ip只支持ipv4", c)
			return
		}
		if cr.Gateway == model.IP {
			res.FailWithMsg("网关ip不能是探针ip", c)
			return
		}
		_, _net, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", model.IP, model.Mask))
		if !_net.Contains(gateway) {
			res.FailWithMsg("网关ip不属于当前子网", c)
			return
		}
	}
	if cr.CanUseHoneyIPRange != "" {
		// 判断ip范围
		ipList, err1 := ip.ParseIPList(cr.CanUseHoneyIPRange)
		if err1 != nil {
			res.FailWithMsg("ip范围格式错误", c)
			return
		}

		// 判断ip在不在子网里面
		for _, s := range ipList {
			if !model.InSubnet(s) {
				res.FailWithMsg(fmt.Sprintf("%s不属于当前子网", s), c)
				return
			}
		}
	}

	err = global.DB.Model(&model).Updates(map[string]any{
		"title":                  cr.Title,
		"gateway":                cr.Gateway,
		"can_use_honey_ip_range": cr.CanUseHoneyIPRange,
	}).Error
	if err != nil {
		res.FailWithMsg("网络信息修改失败", c)
		return
	}
	res.OkWithMsg("网络信息修改成功", c)
	return
}
