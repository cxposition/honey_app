package flags

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"honey_app/apps/honey_server/global"
	"honey_app/apps/honey_server/models"
	"honey_app/apps/honey_server/utils/pwd"
	"os"
	"time"
)

type User struct {
}

type UserInfoRequest struct {
	Role     int8   `json:"role"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (User) Create(value string) {
	var userInfo UserInfoRequest
	if value != "" {
		err := json.Unmarshal([]byte(value), &userInfo)
		if err != nil {
			logrus.Errorf("用户创建失败 %s", err)
			return
		}
	} else {
		fmt.Println("请选择角色： 1 管理员 2 普通用户")
		_, err := fmt.Scanln(&userInfo.Role)
		if err != nil {
			fmt.Println("输入错误", err)
			return
		}
		if !(userInfo.Role == 1 || userInfo.Role == 2) {
			fmt.Println("用户角色输入错误", err)
			return
		}
		fmt.Println("请输入用户名")
		fmt.Scanln(&userInfo.Username)

		fmt.Println("请输入密码")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd())) // 读取用户输入的密码
		if err != nil {
			fmt.Println("读取密码时出错:", err)
			return
		}
		fmt.Println("请再次输入密码")
		rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd())) // 读取用户输入的密码
		if err != nil {
			fmt.Println("读取密码时出错:", err)
			return
		}
		if string(password) != string(rePassword) {
			fmt.Println("两次密码不一致")
			return
		}
		userInfo.Password = string(password)
	}

	var u models.UserModel
	err := global.DB.Take(&u, "username = ?", userInfo.Username).Error
	if err == nil {
		fmt.Println("用户名已存在")
		return
	}

	hashPwd, _ := pwd.GenerateFromPassword(userInfo.Password)
	err = global.DB.Create(&models.UserModel{
		Username: userInfo.Username,
		Password: hashPwd,
		Role:     userInfo.Role,
	}).Error
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
		return
	}
	logrus.Infof("用户创建成功")

}
func (User) List() {
	var userList []models.UserModel
	global.DB.Order("created_at desc").Limit(10).Find(&userList)
	for _, model := range userList {
		fmt.Printf("用户id：%d  用户名：%s 用户角色：%d 创建时间：%s\n",
			model.ID,
			model.Username,
			model.Role,
			model.CreatedAt.Format(time.DateTime),
		)
	}
}
