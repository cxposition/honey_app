package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username      string `json:"username"`
	Role          int8   `json:"role"`
	Password      string `json:"password"`
	LastLoginDate string `json:"lastLoginDate"`
}
