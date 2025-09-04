package models

type HoneyPortModel struct {
	Model
	NodeID       uint         `json:"nodeID"`
	NodeModel    NodeModel    `gorm:"foreignKey:NodeID" json:"-"`
	NetID        uint         `json:"netID"`
	NetModel     NetModel     `gorm:"foreignKey:NetID" json:"-"`
	HoneyIpID    uint         `json:"HoneyIpID"`
	HoneyIpModel HoneyIpModel `gorm:"foreignKey:HoneyIpID" json:"-"`
	ServiceID    uint         `json:"serviceID"` // 服务ID
	Port         int          `json:"port"`      // 服务端口
	DstIP        string       `json:"dstIP"`     // 目标IP
	DstPort      int          `json:"distPort"`  // 目标端口
	Status       int8         `json:"status"`    // 状态
}
