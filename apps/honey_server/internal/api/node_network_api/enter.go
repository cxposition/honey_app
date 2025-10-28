package node_network_api

import "sync"

type NodeNetworkApi struct {
	mutex sync.Mutex
}
