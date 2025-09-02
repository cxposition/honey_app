package models

type UserModel struct {
	Model
	Username      string     `json:"username"`
	Role          int8       `json:"role"` // 1 管理员 2 普通用户
	Password      string     `json:"password"`
	LastLoginDate string     `json:"lastLoginDate"`
	LogList       []LogModel `gorm:"foreignKey:UserID" json:"logList"`
}
