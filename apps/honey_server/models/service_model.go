package models

import "gorm.io/gorm"

type ServiceModel struct {
	gorm.Model
	Title        string `json:"title"`
	Agreement    int8   `json:"agreement"`
	ImageID      uint   `json:"imageID"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Status       int8   `json:"status"`
	HoneyIPCount int    `json:"honeyIPCount"`
	ContainerID  string `json:"containerID"` // 容器ID
}
