package node_network_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

func (n *NodeNetworkApi) EnableView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeNetworkModel
	err := global.DB.Debug().Preload("NodeModel").Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网卡不存在", c)
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	if model.Status == 1 {
		res.FailWithMsg("网卡已启用,请勿重复启用", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 启用网卡
		// 往网络列表中添加一条记录
		var net = models.NetModel{
			NodeID:  model.NodeID,
			Title:   fmt.Sprintf("%s_%s_网络", model.NodeModel.Title, model.Network),
			Network: model.Network,
			IP:      model.IP,
			Mask:    model.Mask,
			Gateway: model.Gateway,
		}

		err = tx.Create(&net).Error
		if err != nil {
			res.FailWithMsg("网络记录添加失败", c)
			return err
		}

		// 往主机表中添加记录
		var host = models.HostModel{
			NodeID:   model.NodeID,
			NetID:    net.ID,
			NetModel: models.NetModel{},
			IP:       net.IP,
		}
		err = tx.Create(&host).Error
		if err != nil {
			res.FailWithMsg("主机记录添加失败", c)
			return err
		}

		// 修改状态
		err = tx.Model(&model).Update("status", 1).Error
		return err
	})
	if err != nil {
		logrus.Errorf("网卡启用失败: %s", err)
		return
	}
	err = global.DB.Create(&models.NodeNetworkModel{}).Error
	res.OkWithMsg("网卡启用成功", c)
}
