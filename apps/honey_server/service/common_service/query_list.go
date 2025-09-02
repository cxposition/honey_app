package common_service

import (
	"fmt"
	"gorm.io/gorm"
	"honey_app/apps/honey_server/core"
	"honey_app/apps/honey_server/models"
)

type Request struct {
	Debug    bool
	Likes    []string
	Where    *gorm.DB
	Preload  []string
	Sort     string
	PageInfo models.PageInfo
}

func QueryList[T any](model T, req Request) (list []T, count int64, err error) {
	db := core.GetDB()
	if req.Debug {
		db = db.Debug()
	}

	// 处理preLoad, 预加载字段
	for _, s := range req.Preload {
		db = db.Preload(s)
	}

	// 针对字段的精确匹配
	db = db.Where(model)

	// 高级查询
	if req.Where != nil {
		db = db.Where(req.Where)
	}

	// 模糊匹配
	if req.PageInfo.Key != "" {
		like := core.GetDB().Where("")
		for _, c := range req.Likes {
			db.Or(fmt.Sprintf("%s like ?", c), fmt.Sprintf("%%%s%%", req.PageInfo.Key))
		}
		db = db.Where(like)
	}

	// 分页
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 10
	}
	if req.PageInfo.Page <= 0 {
		req.PageInfo.Page = 1
	}
	offset := (req.PageInfo.Page - 1) * req.PageInfo.Limit
	err = db.Offset(offset).Limit(req.PageInfo.Limit).Order(req.Sort).Find(&list).Error
	err = db.Count(&count).Error
	return
}
