package models

type HoneyPortModel struct {
	Model
	NodeID    uint   `json:"nodeID"`
	NetID     uint   `json:"netID"`
	HoneyIpID uint   `json:"HoneyIpID"`
	ServiceID uint   `json:"serviceID"` // 服务ID
	Port      int    `json:"port"`      // 服务端口
	DstIP     string `json:"dstIP"`     // 目标IP
	DstPort   int    `json:"distPort"`  // 目标端口
	Status    int8   `json:"status"`    // 状态
}
