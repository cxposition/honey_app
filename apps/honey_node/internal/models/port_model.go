package models

type PortModel struct {
	Model
	LocalAddr  string `gorm:"size:64" json:"localAddr"`
	TargetAddr string `gorm:"size:64" json:"targetAddr"`
}
