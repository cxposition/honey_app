package models

type HoneyIpModel struct {
	Model
	NodeID    uint      `json:"nodeID"`
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"`
	NetID     uint      `json:"netID"`
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Mac       string    `gorm:"size:64" json:"mac"`
	Network   string    `gorm:"size:32" json:"network"` // 网卡
	Status    int8      `json:"status"`                 // 1.创建中 2 运行中 3 失败 4 删除中
}
