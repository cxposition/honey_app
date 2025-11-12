package ip_service

import (
	"honey_node/internal/global"
	"honey_node/internal/models"
)

func IPLoad() {
	var ipList []models.IpModel
	global.DB.Find(&ipList)
	//for _, model := range ipList {
	//	// 判断ip在不在这个机器上，如果在就不创建了
	//}
}
