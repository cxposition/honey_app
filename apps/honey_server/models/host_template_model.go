package models

import "gorm.io/gorm"

type HostTemplateModel struct {
	gorm.Model
	Title    string               `gorm:"size:32" json:"title"`
	PortList HostTemplatePortList `gorm:"serializer:json" json:"portList"` // 主机模版列表
}

type HostTemplatePortList []HostTemplatePort
type HostTemplatePort struct {
	Port      int  `json:"port"`
	ServiceID uint `json:"serviceID"`
}
