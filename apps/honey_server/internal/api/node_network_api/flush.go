package node_network_api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
)

// FlushView
// 下发命令 -> 让节点刷新网卡列表 -> 同步数据库 -> 返回最新网卡状态
func (NodeNetworkApi) FlushView(c *gin.Context) {
	// 1️⃣ 参数解析
	cr := middleware.GetBind[models.IDRequest](c)

	// 2️⃣ 校验节点是否存在且已启动
	var node models.NodeModel
	if err := global.DB.Take(&node, cr.ID).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	if node.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	// 3️⃣ 构造命令请求（让节点刷新网卡）
	taskID := uuid.New().String()
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  taskID,
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	// 4️⃣ 通过 gRPC 下发命令
	resp, err := grpc_service.SendCommand(node.Uid, req, 8*time.Second)
	if err != nil {
		res.FailWithMsg(fmt.Sprintf("命令执行失败: %v", err), c)
		return
	}
	if resp == nil || resp.NetworkFlushOutMessage == nil {
		res.FailWithMsg("节点未返回网卡信息", c)
		return
	}

	newList := resp.NetworkFlushOutMessage.NetworkList
	if len(newList) == 0 {
		res.FailWithMsg("节点返回的网卡列表为空", c)
		return
	}

	// 5️⃣ 过滤 IPv6 地址
	filteredList := make([]*node_rpc.NetworkInfoMessage, 0, len(newList))
	for _, n := range newList {
		if strings.Contains(n.Ip, ":") {
			logrus.Infof("跳过 IPv6 网卡: %s (%s)", n.Network, n.Ip)
			continue
		}
		filteredList = append(filteredList, n)
	}
	if len(filteredList) == 0 {
		res.FailWithMsg("节点仅返回 IPv6 网卡，未发现可用 IPv4", c)
		return
	}

	// 6️⃣ 获取数据库中现有记录
	var oldList []models.NodeNetworkModel
	if err := global.DB.Find(&oldList, "node_id = ?", node.ID).Error; err != nil {
		res.FailWithMsg("查询节点网卡失败", c)
		return
	}

	// 构建 map 方便比对（key: network + ip）
	oldMap := make(map[string]models.NodeNetworkModel)
	for _, old := range oldList {
		key := fmt.Sprintf("%s_%s", old.Network, old.IP)
		oldMap[key] = old
	}

	for _, n := range filteredList {
		key := fmt.Sprintf("%s_%s", n.Network, n.Ip)
		old, exists := oldMap[key]

		if !exists {
			// 🟢 新增
			newRecord := models.NodeNetworkModel{
				NodeID:  node.ID,
				Network: n.Network,
				IP:      n.Ip,
				Mask:    int8(n.Mask),
				Status:  1,
			}
			if err := global.DB.Create(&newRecord).Error; err != nil {
				logrus.Errorf("插入新网卡失败 %s: %v", n.Network, err)
			} else {
				logrus.Infof("新增网卡: %s (%s/%d)", n.Network, n.Ip, n.Mask)
			}
			continue
		}

		if old.Mask != int8(n.Mask) {
			if err := global.DB.Model(&old).Update("mask", int8(n.Mask)).Error; err != nil {
				logrus.Errorf("更新网卡失败 %s: %v", n.Network, err)
			} else {
				logrus.Infof("更新网卡: %s -> mask %d", n.Network, n.Mask)
			}
		}

		// 标记已处理
		delete(oldMap, key)
	}

	// 8️⃣ 删除旧的、节点中不存在的网卡
	for _, obsolete := range oldMap {
		if err := global.DB.Delete(&obsolete).Error; err != nil {
			logrus.Errorf("删除旧网卡失败 %s: %v", obsolete.Network, err)
		} else {
			logrus.Infof("删除旧网卡: %s (%s)", obsolete.Network, obsolete.IP)
		}
	}

	// 9️⃣ 返回更新后的数据
	var updatedList []models.NodeNetworkModel
	if err := global.DB.Find(&updatedList, "node_id = ?", node.ID).Error; err != nil {
		res.FailWithMsg("刷新后查询失败", c)
		return
	}

	res.OkWithData(gin.H{
		"node":      node,
		"network":   updatedList,
		"rawResult": resp.NetworkFlushOutMessage,
	}, c)
}
