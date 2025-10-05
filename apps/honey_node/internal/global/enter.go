package global

import (
	"github.com/sirupsen/logrus"
	"honey_node/internal/config"
)

var (
	Config *config.Config
)

var (
	Version   = "v1.0.1"
	Commit    = "a5f28b47b9"
	BuildTime = "2025-08-21 21:30:05"
)

var (
	Log *logrus.Entry
)
