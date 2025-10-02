package models

type HostTemplateModel struct {
	Model
	Title    string               `gorm:"size:32" json:"title"`
	PortList HostTemplatePortList `gorm:"serializer:json" json:"portList"` // 主机模版列表
}

type HostTemplatePortList []HostTemplatePort
type HostTemplatePort struct {
	Port      int  `json:"port" binding:"min=1,max=65535"`
	ServiceID uint `json:"serviceID" binding:"required"`
}
