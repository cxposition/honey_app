package node_network_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"sync"
)

var mu sync.Mutex

func (n NodeNetworkApi) EnableView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeNetworkModel
	err := global.DB.Debug().Preload("NodeModel").Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("网卡不存在", c)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if model.Status == 1 {
		res.FailWithMsg("网卡已启用,请勿重复启用", c)
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 启用网卡
		// 往网络列表中添加一条记录
		ipRange, err1 := ip.GetUsableIPRange(fmt.Sprintf("%s/%d", model.IP, model.Mask))
		if err1 != nil {
			return err1
		}
		var net = models.NetModel{
			NodeID:             model.NodeID,
			Title:              fmt.Sprintf("%s_%s_网络", model.NodeModel.Title, model.Network),
			Network:            model.Network,
			IP:                 model.IP,
			Mask:               model.Mask,
			Gateway:            model.Gateway,
			CanUseHoneyIPRange: ipRange,
		}

		err = tx.Create(&net).Error
		if err != nil {
			res.FailWithMsg("网络记录添加失败", c)
			return err
		}

		// 往主机表中添加记录
		var host = models.HostModel{
			NodeID: model.NodeID,
			NetID:  net.ID,
			IP:     net.IP,
		}
		err = tx.Create(&host).Error
		if err != nil {
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
	res.OkWithMsg("网卡启用成功", c)
}
