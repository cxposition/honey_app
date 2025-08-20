package main

import (
	"honey_app/apps/honey_server/core"
	"honey_app/apps/honey_server/flags"
	"honey_app/apps/honey_server/global"
)

func main() {
	global.DB = core.InitDB()
	flags.Run()
}
