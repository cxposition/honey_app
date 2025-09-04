package common_service

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"honey_app/apps/honey_server/internal/global"
)

type RemoveRequest struct {
	Debug    bool
	Where    *gorm.DB
	IDList   []uint
	Log      *logrus.Entry
	Msg      string
	Unscoped bool
}

func Remove[T any](model T, req RemoveRequest) (successCount int64, err error) {
	db := global.DB
	if req.Debug {
		db = db.Debug()
	}
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	if req.Unscoped {
		req.Log.Infof("真删除")
		db = db.Unscoped()
	}

	db = db.Where(model)
	if len(req.IDList) > 0 {
		req.Log.Infof("删除 %s idList %v", req.Msg, req.IDList)
		db = db.Where("id in ?", req.IDList)
	}

	var list []T
	db.Find(&list)
	if len(list) <= 0 {
		req.Log.Infof("没查到")
		return
	}

	result := db.Delete(&list)
	if result.Error != nil {
		req.Log.Infof("删除失败 %s", result.Error)
		return
	}
	successCount = result.RowsAffected
	req.Log.Infof("删除 %s 成功 , 成功%v个", req.Msg, successCount)
	return
}
