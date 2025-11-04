package net_api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
	"time"
)

func (NetApi) ScanView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NetModel
	if err := global.DB.Preload("NodeModel").Take(&model, cr.ID).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	if model.NodeModel.Status != 1 {
		res.FailWithMsg("节点未启动", c)
		return
	}

	// 3️⃣ 构造命令请求（让节点刷新网卡）
	taskID := uuid.New().String()
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  taskID,
		NetScanInMessage: &node_rpc.NetScanInMessage{
			Network:      model.Network,
			IpRange:      model.CanUseHoneyIPRange,
			FilterIPList: []string{},
			NetID:        uint32(model.ID),
		},
	}

	nodeID := model.NodeModel.Uid
	val, ok := grpc_service.NodeCommandMap.Load(nodeID)
	if !ok {
		logrus.Errorf("节点 %s 未找到", nodeID)
		//cancelFunc()
		return
	}
	cmd := val.(*grpc_service.Command)

	respChan := make(chan *node_rpc.CmdResponse, 1)
	cmd.ResMap.Store(req.TaskID, respChan)
	defer cmd.ResMap.Delete(req.TaskID)

	select {
	case cmd.ReqChan <- req:
		// 成功发送
		logrus.Infof("节点 %s 的 task %s 已发送", nodeID, req.TaskID)
	case <-time.After(3 * time.Second):
		logrus.Errorf("发送命令到节点 %s 超时", nodeID)
	}

	// ✅ 添加超时机制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var netScanMsg []*node_rpc.NetScanOutMessage
label:
	for {
		select {
		case resp := <-respChan:
			logrus.Debugf("节点 %s 的 task %s 的结果为 %+v", nodeID, req.TaskID, resp)
			message := resp.NetScanOutMessage
			logrus.Infof("节点信息为:%+v", message)
			if message.ErrMsg != "" {
				res.FailWithMsg("扫描错误"+message.ErrMsg, c)
				break label
			}
			if message.End {
				break label
			}
			netScanMsg = append(netScanMsg, message)
		case <-ctx.Done():
			res.FailWithMsg("获取响应超时", c)
			return
		default:
		}
	}

	// 当前的主机列表
	var hostList []models.HostModel
	global.DB.Find(&hostList, "net_id = ?", cr.ID)

	// 算出新增的主机,删除的主机

	// -----------------------------
	// ✅ 算出新增的主机和删除的主机
	// -----------------------------

	// 1️⃣ 把数据库中已有主机转成 map[ip]HostModel
	dbHosts := make(map[string]models.HostModel)
	for _, h := range hostList {
		dbHosts[h.IP] = h
	}

	// 2️⃣ 把扫描结果转成 map[ip]*node_rpc.NetScanOutMessage
	scanHosts := make(map[string]*node_rpc.NetScanOutMessage)
	for _, msg := range netScanMsg {
		scanHosts[msg.Ip] = msg
	}

	// 3️⃣ 计算新增主机：在扫描结果中但不在数据库中
	var newHosts []models.HostModel
	for ip, msg := range scanHosts {
		if _, exists := dbHosts[ip]; !exists {
			newHost := models.HostModel{
				NodeID: model.NodeID,
				NetID:  model.ID,
				IP:     msg.Ip,
				Mac:    msg.Mac,
				Manuf:  msg.Manuf,
			}
			newHosts = append(newHosts, newHost)
		}
	}

	// 4️⃣ 计算删除主机：在数据库中但不在扫描结果中
	var delIPs []string
	for ip := range dbHosts {
		if _, exists := scanHosts[ip]; !exists {
			delIPs = append(delIPs, ip)
		}
	}

	// 5️⃣ 执行数据库更新
	tx := global.DB.Begin()

	if len(newHosts) > 0 {
		if err := tx.Create(&newHosts).Error; err != nil {
			tx.Rollback()
			res.FailWithMsg("新增主机保存失败: "+err.Error(), c)
			return
		}
		logrus.Infof("网络 %d 扫描新增主机 %d 个", model.ID, len(newHosts))
	}

	if len(delIPs) > 0 {
		if err := tx.Where("net_id = ? AND ip IN ?", model.ID, delIPs).Delete(&models.HostModel{}).Error; err != nil {
			tx.Rollback()
			res.FailWithMsg("删除主机失败: "+err.Error(), c)
			return
		}
		logrus.Infof("网络 %d 扫描删除主机 %d 个", model.ID, len(delIPs))
	}

	tx.Commit()

	res.OkWithMsg("扫描成功", c)
}
