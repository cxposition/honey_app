package core

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"honey_node/internal/config"
	"honey_node/internal/flags"
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
		logrus.Fatalf("配置文件配置错误 %s", err)
		return nil
	}
	SetDefault(&c)
	return &c
}

func SetDefault(c *config.Config) {
	if c.System.Uid == "" {
		c.System.Uid = uuid.NewString()
		SetConfig(c)
	}
}

func SetConfig(c *config.Config) {
	byteData, err := yaml.Marshal(c)
	if err != nil {
		logrus.Errorf("配置序列化失败 %s", err)
		return
	}

	err = os.WriteFile(flags.Options.File, byteData, 0666)
	if err != nil {
		logrus.Errorf("配置文件写入错误 %s", err)
		return
	}
	logrus.Infof("%s 配置文件更新成功", flags.Options.File)
}
