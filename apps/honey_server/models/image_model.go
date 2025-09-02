package models

type ImageModel struct {
	Model
	ImageName     string `json:"imageName"`
	Title         string `json:"title"`
	Port          int    `json:"port"`
	DockerImageID string `gorm:"size:32" json:"dockerImageID"`
	Tag           string `json:"tag"`
	Agreement     int8   `json:"agreement"`
	ImagePath     string `json:"imagePath"` // 镜像文件
	Status        int8   `json:"status"`
	Logo          string `json:"logo"` // 镜像logo
	Desc          string `json:"desc"` // 镜像描述
}
