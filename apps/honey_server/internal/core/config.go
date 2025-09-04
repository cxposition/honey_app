package core

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"honey_server/internal/config"
	"honey_server/internal/flags"
	"os"
)

func ReadConfig() *config.Config {
	byteData, err := os.ReadFile(flags.Options.File)
	if err != nil {
		logrus.Fatalf("配置文件读取错误 %s", err)
		return nil
	}

	var c config.Config
	err = yaml.Unmarshal(byteData, &c)
	if err != nil {
		logrus.Fatal("配置文件配置错误 %s", err)
		return nil
	}
	return &c
}
