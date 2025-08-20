package models

import "gorm.io/gorm"

type MatrixTemplateModel struct {
	gorm.Model
	Title            string           `gorm:"size:32" json:"title"`
	HostTemplateList HostTemplateList `gorm:"serializer:json" json:"hostTemplateList"`
}

type HostTemplateList []HostTemplateInfo
type HostTemplateInfo struct {
	HostTemplateID uint `json:"hostTemplateID"`
	Weight         int  `json:"weight"`
}
