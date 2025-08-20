package models

import "gorm.io/gorm"

type ServiceModel struct {
	gorm.Model
	Title         string     `json:"title"`
	Agreement     int8       `json:"agreement"`
	ImageID       uint       `json:"imageID"`
	ImageModel    ImageModel `gorm:"foreignKey:ImageID" json:"-"`
	IP            string     `json:"ip"`
	Port          int        `json:"port"`
	Status        int8       `json:"status"`
	HoneyIPCount  int        `json:"honeyIPCount"`
	ContainerID   string     `json:"containerID"`                  // 容器ID
	ContainerName string     `gorm:"size:32" json:"containerName"` // 容器名
}
