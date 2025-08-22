package flags

import (
	"flag"
	"github.com/sirupsen/logrus"
	"honey_app/apps/honey_server/global"
	"os"
)

type FlagOptions struct {
	File    string
	Version bool
	DB      bool
}

var Options FlagOptions

func init() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "v", false, "打印当前版本")
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.Parse()
}

func Run() {
	if Options.DB {
		Migrate()
		os.Exit(0)
	}
	if Options.Version {
		logrus.Infof("当前版本: %s commit: %s buildTime: %s",
			global.Version,
			global.Commit,
			global.BuildTime,
		)
		os.Exit(0)
	}
}
