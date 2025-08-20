package models

import "gorm.io/gorm"

type HostTemplateModel struct {
	gorm.Model
	Title            string           `gorm:"32" json:"title"`
	HostTemplateList HostTemplateList `json:"hostTemplateList"` // 主机模版列表
}

type HostTemplatePortList []HostTemplatePort
type HostTemplatePort struct {
	HostTemplateID uint `json:"hostTemplateID"`
	Weight         int  `json:"weight"`
}
