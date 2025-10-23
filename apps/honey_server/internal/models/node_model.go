package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NodeModel 节点表
type NodeModel struct {
	Model
	Title        string `json:"title"` // 节点名称
	Uid          string `gorm:"size:64" json:"uid"`
	IP           string `json:"IP"`
	Mac          string `gorm:"size:64" json:"mac"`
	Status       int8   `json:"status"`       // 节点状态
	NetCount     int    `json:"netCount"`     // 网络连接数目
	HoneyIPCount int    `json:"honeyIPCount"` // 诱捕IP数
	// 查资料说必须要加serializer，不然会报错，gorm默认不支持自定义的数据类型序列化和反序列化
	Resource   NodeResource   `gorm:"serializer:json" json:"resource"`
	SystemInfo NodeSystemInfo `gorm:"serializer:json" json:"systemInfo"`
}

func (n *NodeModel) BeforeDelete(tx *gorm.DB) error {
	// 诱捕转发
	var list []HoneyPortModel
	err := tx.Find(&list, "node_id = ?", n.ID).Delete(&list).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联诱捕转发 %d", len(list))

	// 诱捕ip
	var ipList []HoneyIpModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&ipList).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联诱捕ip %d", len(ipList))

	// 节点网络
	var netList []NetModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&netList).Error
	if err != nil {
		return err
	}
	logrus.Infof("节点网络 %d", len(netList))

	// 节点网卡
	var networkList []NodeNetworkModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&networkList).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联节点网卡 %d", len(networkList))

	// 节点
	//var list []HoneyPortModel
	//tx.Find(&list, "node_id = ?", n.ID).Delete(&list)
	//logrus.Infof("关联诱捕节点 %d", len(list))

	return nil
}

type NodeResource struct {
	CpuCount              int     `json:"cpuCount"`
	CpuUseRate            float64 `json:"cpuUseRate"`
	MemTotal              int64   `json:"memTotal"`
	MemUseRate            float64 `json:"memUseRate"`
	DiskTotal             int64   `json:"diskTotal"`
	DiskUseRate           float64 `json:"diskUseRate"`
	NodePath              string  `json:"nodePath"`
	NodeResourceOccupancy int64   `json:"nodeResourceOccupancy"`
}

type NodeSystemInfo struct {
	HostName            string `json:"hostName"`
	DistributionVersion string `json:"distributionVersion"` // 发行版本
	CoreVersion         string `json:"coreVersion"`         // 内核版本
	SystemType          string `json:"systemType"`          // 系统类型
	StartTime           string `json:"startTime"`           // 启动时间
	NodeVersion         string `json:"nodeVersion"`
	NodeCommit          string `json:"nodeCommit"`
}

type IDRequest struct {
	ID []uint `json:"id" form:"id" uri:"id"`
}
