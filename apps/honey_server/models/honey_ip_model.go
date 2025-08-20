package models

import "gorm.io/gorm"

type HoneyIpModel struct {
	gorm.Model
	NodeID    uint      `json:"nodeID"`
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"`
	NetID     uint      `json:"netID"`
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`
	IP        string    `gorm:"32" json:"ip"`
	Mac       string    `gorm:"64" json:"mac"`
	Network   string    `gorm:"32" json:"network"` // 网卡
	Status    int8      `json:"status"`
}
