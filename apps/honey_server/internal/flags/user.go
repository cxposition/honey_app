package flags

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"honey_app/apps/honey_server/internal/global"
	"honey_app/apps/honey_server/internal/models"
	user_service2 "honey_app/apps/honey_server/internal/service/user_service"
	"os"
	"time"
)

type User struct {
}

func (User) Create(value string) {
	var userInfo user_service2.UserCreateRequest
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

	us := user_service2.NewUserService(global.Log)
	_, err := us.Create(userInfo)
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
	}

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
