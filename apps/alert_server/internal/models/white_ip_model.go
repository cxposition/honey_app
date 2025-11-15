package models

type WhiteIPModel struct {
	Model
	IP     string `gorm:"size:32" json:"ip"`
	Notice string `gorm:"size:64" json:"notice"`
}
