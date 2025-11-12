package models

type IpModel struct {
	Model
	Ip       string `json:"ip"`
	Mask     int8   `json:"mask"`
	LinkName string `json:"linkName"` // 自己的接口名称
	Network  string `json:"network"`  // 基于哪个网卡创建
	Mac      string `json:"mac"`
}
