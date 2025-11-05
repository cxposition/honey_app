package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net"
)

// NetModel 网络表
type NetModel struct {
	Model
	NodeID             uint      `json:"nodeID"`
	NodeModel          NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 这个foreignKey要写结构体字段名
	Title              string    `gorm:"size:32" json:"title"`
	Network            string    `gorm:"size:32" json:"network"` // 网卡
	IP                 string    `gorm:"size:32" json:"ip"`      // 探针ip
	Mask               int8      `json:"mask"`                   // 子网掩码 8-32
	Gateway            string    `gorm:"size:32" json:"gateway"`
	HostCount          int       `json:"hostCount"`                          // 存活资产
	HoneyIpCount       int       `json:"honeyIpCount"`                       // 诱捕ip
	ScanStatus         int8      `json:"scanStatus"`                         // 扫描状态 0表示待扫描 1扫描完成 2扫描中
	ScanProgress       float64   `json:"scanProgress"`                       // 扫描进度
	CanUseHoneyIPRange string    `gorm:"size:256" json:"canUseHoneyIpRange"` // 能够使用的诱捕ip范围
}

func (model NetModel) Subnet() string {
	return fmt.Sprintf("%s/%d", model.IP, model.Mask)
}

func (model NetModel) InSubnet(ip string) bool {
	_, _net, _ := net.ParseCIDR(model.Subnet())
	logrus.Infof("_net:%+v", _net)
	return _net.Contains(net.ParseIP(ip))
}

func (model NetModel) BeforeDelete(tx *gorm.DB) error {
	// 是否有诱捕ip
	var count int64
	tx.Model(&HoneyIpModel{}).Where("net_id = ?", model.ID).Count(&count)
	if count > 0 {
		return errors.New("存在诱捕ip，不能删除网络")
	}

	// 将启用的网卡状态归位
	var nodeNet NodeNetworkModel
	err := tx.Take(&nodeNet, "net_id = ? and network = ?", model.ID, model.Network).Error
	if err != nil {
		return nil
	}

	var hostList []HostModel
	tx.Find(&hostList, "net_id = ?", model.ID).Delete(&hostList)
	logrus.Infof("关联删除主机%d个", len(hostList))

	// 修改状态
	err = tx.Model(&nodeNet).Update("status", 2).Error
	return nil
}
