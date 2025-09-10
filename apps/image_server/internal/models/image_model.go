package models

type ImageModel struct {
	Model
	ImageName     string `gorm:"size:64" json:"imageName"`
	Title         string `gorm:"size:64" json:"title"`
	Port          int    `json:"port"`
	DockerImageID string `gorm:"size:32" json:"dockerImageID"`
	Tag           string `gorm:"size:32" json:"tag"`
	Agreement     int8   `json:"agreement"`                 // 协议
	ImagePath     string `gorm:"size:256" json:"imagePath"` // 镜像文件
	Status        int8   `json:"status"`
	Logo          string `gorm:"size:256" json:"logo"` // 镜像logo
	Desc          string `gorm:"size:512" json:"desc"` // 镜像描述
}
