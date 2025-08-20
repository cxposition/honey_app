package models

import "gorm.io/gorm"

type ImageModel struct {
	gorm.Model
	ImageName  string `json:"imageName"`
	Title      string `json:"title"`
	Port       int    `json:"port"`
	ImageID    string `json:"imageID"`
	ImageModel `gorm:"foreignKey:ImageID" json:"-"`
	Tag        string `json:"tag"`
	Agreement  int8   `json:"agreement"`
	ImagePath  string `json:"imagePath"` // 镜像文件
	Status     int8   `json:"status"`
	Logo       string `json:"logo"` // 镜像logo
	Desc       string `json:"desc"` // 镜像描述
}
