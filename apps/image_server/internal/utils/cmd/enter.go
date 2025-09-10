package cmd

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Cmd(command string) (err error) {
	// 创建一个Cmd结构体
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)
	// 设置输出
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	// 执行命令
	err = cmd.Run()
	if err != nil {
		return err
	}
	logrus.Infof("命令输出 %s", stdout.String())
	return nil
}
