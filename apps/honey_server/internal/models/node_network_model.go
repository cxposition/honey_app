package models

// NodeNetworkModel 节点网络表
type NodeNetworkModel struct {
	Model
	NodeID    uint      `json:"nodeID"`
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 这个foreignKey要写结构体字段名
	Network   string    `json:"network"`                    // 网卡
	IP        string    `json:"ip"`                         // 探针ip
	Mask      int8      `json:"mask"`                       // 子网掩码 8-32
	Gateway   string    `json:"gateway"`
	Status    int8      `json:"status"` // 是否启用
	NetworkID uint      `json:"networkID"`
}
