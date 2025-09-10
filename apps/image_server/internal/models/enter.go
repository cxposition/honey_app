package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type PageInfo struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Key   string `form:"key"`
}

type IDListRequest struct {
	IDList []uint `json:"idList"`
}

type IDRequest struct {
	ID []uint `json:"id" form:"id" uri:"id"`
}

type RemoveRequest struct {
	Debug    bool
	IDList   []uint
	Log      *logrus.Entry
	Msg      string
	Unscoped bool
	Where    *gorm.DB
}
