package models

import "gorm.io/gorm"

// NetModel 网络表
type NetModel struct {
	gorm.Model
	NodeID             uint      `json:"nodeID"`
	NodeModel          NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 这个foreignKey要写结构体字段名
	Title              string    `gorm:"32" json:"title"`
	Network            string    `gorm:"32" json:"network"` // 网卡
	IP                 string    `gorm:"32" json:"ip"`      // 探针ip
	Mask               int8      `json:"mask"`              // 子网掩码 8-32
	Gateway            string    `gorm:"32" json:"gateway"`
	HostCount          int       `json:"hostCount"`                     // 存活资产
	HoneyIpCount       int       `json:"honeyIpCount"`                  // 诱捕ip
	ScanStatus         int8      `json:"scanStatus"`                    // 扫描状态
	ScanProgress       float64   `json:"scanProgress"`                  // 扫描进度
	CanUseHoneyIPRange string    `gorm:"256" json:"canUseHoneyIpRange"` // 能够使用的诱捕ip范围
}
