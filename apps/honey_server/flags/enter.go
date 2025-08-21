package flags

import (
	"flag"
	"fmt"
	"os"
)

type FlagOptions struct {
	File    string
	Version bool
	DB      bool
}

var Options FlagOptions

func init() {
	dir, _ := os.Getwd()
	fmt.Println(dir)
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
}
