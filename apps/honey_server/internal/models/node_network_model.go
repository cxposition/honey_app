package models

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NodeNetworkModel 节点网络表
type NodeNetworkModel struct {
	Model
	NodeID    uint      `json:"nodeID"`
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 这个foreignKey要写结构体字段名
	Network   string    `json:"network"`                    // 网卡
	IP        string    `json:"ip"`                         // 探针ip
	Mask      int8      `json:"mask"`                       // 子网掩码 8-32
	Gateway   string    `json:"gateway"`
	Status    int8      `json:"status"` // 是否启用 1 表示启用 2表示未启用
	NetworkID uint      `json:"networkID"`
}

func (n *NodeNetworkModel) BeforeDelete(tx *gorm.DB) error {
	// 先找有没有网络
	if n.Status == 2 {
		return nil
	}
	var net NetModel
	err := tx.Take(&net, "node_id = ? and network = ?", n.NodeID, n.Network).Error
	if err != nil {
		// 未启用
		return nil
	}

	// 判断有没有诱捕ip
	var count int64
	err = tx.Model(&HoneyIpModel{}).Where("network_id = ?", net.ID).Count(&count).Error
	if count > 0 {
		return errors.New("此网卡的网络存在诱捕ip，不可删除")
	}

	// 关联删除网络表和主机表
	var hostList []HostModel
	tx.Find(&hostList, "net_id = ?", net.ID).Delete(&hostList)
	tx.Delete(&net)
	logrus.Infof("关联删除主机表%d 和网络表 %s", len(hostList), net.Title)
	return nil
}
