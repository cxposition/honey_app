package models

import "gorm.io/gorm"

// HostModel 主机表
type HostModel struct {
	gorm.Model
	NodeID    uint      `json:"nodeID"`
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 这个foreignKey要写结构体字段名
	NetID     uint      `json:"netID"`
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Mac       string    `gorm:"size:64" json:"mac"`
	Manuf     string    `gorm:"size:64" json:"manuf"` // 厂商信息
}
