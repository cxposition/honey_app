package models

import "gorm.io/gorm"

// NodeModel 节点表
type NodeModel struct {
	gorm.Model
	Title        string `json:"title"` // 节点名称
	IP           string `json:"IP"`
	Status       int8   `json:"status"`       // 节点状态
	NetCount     int    `json:"netCount"`     // 网络连接数目
	HoneyIPCount int    `json:"honeyIPCount"` // 诱捕IP数
	// 查资料说必须要加serializer，不然会报错，gorm默认不支持自定义的数据类型序列化和反序列化
	Resource   NodeResource   `gorm:"serializer:json" json:"resource"`
	SystemInfo NodeSystemInfo `gorm:"serializer:json" json:"systemInfo"`
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
}
