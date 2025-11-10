package global

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"honey_node/internal/config"
	"honey_node/internal/rpc/node_rpc"
)

var (
	Version   = "v1.0.1"
	Commit    = "a5f28b47b9"
	BuildTime = "2025-08-21 21:30:05"
)

var (
	Log        *logrus.Entry
	GrpcClient node_rpc.NodeServiceClient
	Config     *config.Config
	Conn       *amqp.Connection
)
