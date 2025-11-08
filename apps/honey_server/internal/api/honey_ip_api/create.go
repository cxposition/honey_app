package honey_ip_api

import (
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/service/mq_service"
	"honey_server/internal/utils"
	"honey_server/internal/utils/res"
)

type CreateRequest struct {
	NetID uint   `json:"netID" binding:"required"`
	IP    string `json:"ip" binding:"required"`
}

func (HoneyIPApi) CreateView(c *gin.Context) {
	cr := middleware.GetBind[CreateRequest](c)

	// 判断网络是否存在
	var netModel models.NetModel
	err := global.DB.Preload("NodeModel").Take(&netModel, cr.NetID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 判断ip是否在子网 ip range里面
	ipRange, err := netModel.IpRange()
	if !utils.Inlist(ipRange, cr.IP) {
		res.FailWithMsg("当前ip不存在可部署ip列表里面", c)
		return
	}

	// 判断ip是不是主机ip
	var hostModel models.HostModel
	err = global.DB.Take(&hostModel, "net_id = ? and ip = ?", cr.NetID, cr.IP).Error
	if err == nil {
		res.FailWithMsg("当前ip是主机ip", c)
		return
	}

	// 判断ip是不是已经部署过了
	var honeyIPModel models.HoneyIpModel
	err = global.DB.Take(&honeyIPModel, "net_id = ? and ip = ?", cr.NetID, cr.IP).Error
	if err == nil {
		res.FailWithMsg("当前ip已使用", c)
		return
	}

	// 判断节点是否在线
	if netModel.NodeModel.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	_, ok := grpc_service.GetNodeCommand(netModel.NodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	var model = models.HoneyIpModel{
		NodeID: netModel.NodeID,
		NetID:  netModel.ID,
		IP:     cr.IP,
		Status: 1,
	}

	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("创建诱捕ip失败", c)
		return
	}

	// 下发消息
	mq_service.SendCeateIPMsg(netModel.NodeModel.Uid, mq_service.CreateIPRequest{
		HoneyIPID: model.ID,
		IP:        model.IP,
		Mask:      netModel.Mask,
		Network:   netModel.Network,
	})

	res.OkWithData(model.ID, c)
}
