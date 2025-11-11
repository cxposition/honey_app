package models

type PortModel struct {
	Model
	IP       string `gorm:"size:32" json:"ip"`
	Port     int    `json:"port"`
	DestIP   string `gorm:"size:32" json:"destIP"`
	DestPort int    `json:"destPort"`
}
