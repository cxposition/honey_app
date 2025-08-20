package models

import "gorm.io/gorm"

type LogModel struct {
	gorm.Model
	Type        int8   `json:"type"`
	IP          string `json:"ip"`
	Addr        string `json:"addr"`
	UserID      string `json:"userID"`
	Username    string `json:"username"`
	Pwd         string `json:"pwd"`
	LoginStatus bool   `json:"loginStatus"`
	Title       string `json:"title"`
	Level       int8   `json:"level"`
	Content     string `json:"content"`
	ServiceName string `json:"serviceName"`
}
