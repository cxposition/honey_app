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
	Menu    string
	Type    string
	Value   string
}

var Options FlagOptions

func init() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "vv", false, "打印当前版本")
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.StringVar(&Options.Menu, "m", "", "菜单 user")
	flag.StringVar(&Options.Type, "t", "", "类型 create list")
	flag.StringVar(&Options.Value, "v", "", "值")
	flag.Parse()
}

func Run() {
	if Options.DB {
		Migrate()
		os.Exit(0)
	}
	if Options.Version {
		logrus.Infof("当前版本: %s commit: %s buildTime: %s",
			global.Version, global.Commit, global.BuildTime,
		)
		os.Exit(0)
	}
	switch Options.Menu {
	case "user":
		var user User
		switch Options.Type {
		case "create":
			user.Create()
			os.Exit(0)
		case "list":
			user.List()
			os.Exit(0)
		default:
			logrus.Fatalf("用户子菜单项不正确.")
		}
	case "":
	default:
		logrus.Fatalf("菜单项不正确")
	}

}
