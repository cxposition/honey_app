package user_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/utils/pwd"
)

type UserCreateRequest struct {
	Password string `json:"password" binding:"required"`
	Role     int8   `json:"role" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (u *UserService) Create(req UserCreateRequest) (user models.UserModel, err error) {
	err = global.DB.Take(&user, "username = ?", req.Username).Error
	if err == nil {
		fmt.Println("用户名已存在")
		return
	}
	hashPwd, _ := pwd.GenerateFromPassword(req.Password)
	user = models.UserModel{
		Username: req.Username,
		Password: hashPwd,
		Role:     req.Role,
	}
	err = global.DB.Create(&user).Error
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
		return
	}
	u.log.Infof("用户创建成功")
	return
}
