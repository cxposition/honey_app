package cron_service

import (
	"github.com/robfig/cron/v3"
	"time"
)

func Run() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	crontab.AddFunc("*/5 * * * * *", Resource) // 每5s执行 一次
	crontab.Start()
}
