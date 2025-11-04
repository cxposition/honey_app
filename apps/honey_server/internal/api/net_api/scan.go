package net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
	"sync"
	"time"
)

var mux sync.Mutex

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

	// 需要将诱捕ip过滤
	var filterIPList []string
	global.DB.Model(&models.HoneyIpModel{}).Where("net_id = ?", cr.ID).Select("ip").Scan(&filterIPList)
	fmt.Println("过滤的ip列表", filterIPList)

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
		res.FailWithMsg("节点未注册", c)
		return
	}
	cmd := val.(*grpc_service.Command)

	respChan := make(chan *node_rpc.CmdResponse, 100) // 缓冲更大些
	cmd.ResMap.Store(taskID, respChan)

	mux.Lock()
	if model.ScanStatus == 2 {
		res.FailWithMsg("当前子网正在扫描中", c)
		mux.Unlock()
		return
	}

	// 修改状态为扫描中
	global.DB.Model(&model).Update("scan_status", 2)
	mux.Unlock()

	// 异步处理扫描结果
	go handleScanResult(nodeID, model, taskID, respChan)

	// 异步发送请求（防止阻塞）
	go func() {
		select {
		case cmd.ReqChan <- req:
			logrus.Infof("节点 %s 的 task %s 已下发", nodeID, taskID)
		case <-time.After(3 * time.Second):
			logrus.Errorf("发送命令到节点 %s 超时", nodeID)
		}
	}()

	// 立即响应前端
	res.OkWithData(gin.H{
		"task_id": taskID,
		"msg":     "扫描任务已下发，正在后台执行",
	}, c)
}

// handleScanResult 异步接收扫描结果并更新数据库
func handleScanResult(nodeID string, model models.NetModel, taskID string, respChan chan *node_rpc.CmdResponse) {
	defer func() {
		if val, ok := grpc_service.NodeCommandMap.Load(nodeID); ok {
			val.(*grpc_service.Command).ResMap.Delete(taskID)
		}
		close(respChan)
	}()

	defer func() {
		// 函数走完，将状态修改为扫描完成
		global.DB.Model(&model).Update("scan_status", 1)
	}()

	logrus.Infof("[异步扫描] 网络 %d（节点 %s）任务 %s 开始接收结果", model.ID, nodeID, taskID)

	var netScanMsg []*node_rpc.NetScanOutMessage
	timeout := time.After(30 * time.Second) // 最大等待时间 30s

	for {
		select {
		case resp := <-respChan:
			if resp == nil || resp.NetScanOutMessage == nil {
				continue
			}
			msg := resp.NetScanOutMessage
			if msg.ErrMsg != "" {
				logrus.Errorf("扫描错误: %s", msg.ErrMsg)
				return
			}
			if msg.End {
				logrus.Infof("[异步扫描] 网络 %d 扫描完成，共收到 %d 条记录", model.ID, len(netScanMsg))
				updateHosts(model, netScanMsg)
				return
			}
			netScanMsg = append(netScanMsg, msg)

		case <-timeout:
			logrus.Warnf("[异步扫描] 网络 %d 超时，收到 %d 条记录", model.ID, len(netScanMsg))
			updateHosts(model, netScanMsg)
			return
		}
	}
}

// updateHosts 比对数据库，计算新增和删除主机
func updateHosts(model models.NetModel, netScanMsg []*node_rpc.NetScanOutMessage) {
	var hostList []models.HostModel
	global.DB.Find(&hostList, "net_id = ?", model.ID)

	dbHosts := make(map[string]models.HostModel)
	for _, h := range hostList {
		dbHosts[h.IP] = h
	}

	scanHosts := make(map[string]*node_rpc.NetScanOutMessage)
	for _, msg := range netScanMsg {
		scanHosts[msg.Ip] = msg
	}

	var newHosts []models.HostModel
	for ip, msg := range scanHosts {
		if _, exists := dbHosts[ip]; !exists {
			newHosts = append(newHosts, models.HostModel{
				NodeID: model.NodeID,
				NetID:  model.ID,
				IP:     msg.Ip,
				Mac:    msg.Mac,
				Manuf:  msg.Manuf,
			})
		}
	}

	var delIPs []string
	for ip := range dbHosts {
		if _, exists := scanHosts[ip]; !exists {
			delIPs = append(delIPs, ip)
		}
	}

	tx := global.DB.Begin()
	if len(newHosts) > 0 {
		if err := tx.Create(&newHosts).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("新增主机失败: %v", err)
			return
		}
		logrus.Infof("[异步扫描] 网络 %d 新增主机 %d 个", model.ID, len(newHosts))
	}

	if len(delIPs) > 0 {
		if err := tx.Where("net_id = ? AND ip IN ?", model.ID, delIPs).Delete(&models.HostModel{}).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("删除主机失败: %v", err)
			return
		}
		logrus.Infof("[异步扫描] 网络 %d 删除主机 %d 个", model.ID, len(delIPs))
	}

	tx.Commit()
}
