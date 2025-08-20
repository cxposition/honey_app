package models

import "gorm.io/gorm"

type MatrixTemplateModel struct {
	gorm.Model
	Title            string           `gorm:"32" json:"title"`
	HostTemplateList HostTemplateList `gorm:"serializer:json" json:"hostTemplateList"`
}

type HostTemplateList []HostTemplateInfo
type HostTemplateInfo struct {
	Port      int  `json:"port"`
	ServiceID uint `json:"serviceID"`
}
