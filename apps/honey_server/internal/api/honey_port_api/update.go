package honey_port_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/service/mq_service"
	"honey_server/internal/utils/res"
)

type UpdateRequest struct {
	HoneyIPID uint       `json:"honeyIpID" binding:"required"`
	PortList  []PortType `json:"portList" binding:"dive,required"`
}

type PortType struct {
	Port      int  `json:"port" binding:"required,min=1,max=65535"`
	ServiceID uint `json:"serviceID" binding:"required"`
}

func (HoneyPortApi) UpdateView(c *gin.Context) {
	cr := middleware.GetBind[UpdateRequest](c)

	var honeyIPModel models.HoneyIpModel
	err := global.DB.Preload("NodeModel").Take(&honeyIPModel, cr.HoneyIPID).Error
	if err != nil {
		res.FailWithMsg("不存在的诱捕ip", c)
		return
	}

	nodeModel := honeyIPModel.NodeModel

	// 判断节点是否在线
	if nodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 使用封装的获取节点函数
	_, ok := grpc_service.GetNodeCommand(nodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 找之前的端口信息
	var honeyPortList []models.HoneyPortModel
	global.DB.Find(&honeyPortList, "honey_ip_id = ?", cr.HoneyIPID)

	// 服务id有效性的判断
	// 端口不能重复
	var portMap = map[int]struct{}{}
	var serviceIDList []uint
	for _, portType := range cr.PortList {
		serviceIDList = append(serviceIDList, portType.ServiceID)
		portMap[portType.Port] = struct{}{}
	}

	if len(portMap) != len(cr.PortList) {
		res.FailWithMsg("端口重复", c)
		return
	}

	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	var serviceMap = map[uint]models.ServiceModel{}
	for _, model := range serviceList {
		serviceMap[model.ID] = model
	}

	// 算出新增的，删除的
	// 将现有端口转换为map以便快速查找
	existingPorts := make(map[int]models.HoneyPortModel)
	for _, port := range honeyPortList {
		existingPorts[port.Port] = port
	}

	// 计算新增的端口
	var newPorts []models.HoneyPortModel
	for _, reqPort := range cr.PortList {
		service, ok := serviceMap[reqPort.ServiceID]
		if !ok {
			res.FailWithMsg(fmt.Sprintf("服务%d不存在", reqPort.ServiceID), c)
			return
		}

		if _, exists := existingPorts[reqPort.Port]; !exists {
			newPorts = append(newPorts, models.HoneyPortModel{
				HoneyIpID: cr.HoneyIPID,
				Port:      reqPort.Port,
				ServiceID: reqPort.ServiceID,
				DstIP:     service.IP,
				DstPort:   service.Port,
				Status:    1, // 假设1为启用状态
			})
		}
	}

	// 计算需要删除的端口
	var portsToDelete []models.HoneyPortModel
	for port, model := range existingPorts {
		found := false
		for _, reqPort := range cr.PortList {
			if reqPort.Port == port {
				found = true
				break
			}
		}
		if !found {
			portsToDelete = append(portsToDelete, model)
		}
	}

	// 执行数据库操作
	tx := global.DB.Begin()
	if tx.Error != nil {
		res.FailWithMsg("更新端口信息失败", c)
		return
	}

	// 删除端口
	for _, port := range portsToDelete {
		if err := tx.Delete(&port).Error; err != nil {
			tx.Rollback()
			res.FailWithMsg("更新端口信息失败", c)
			return
		}
	}

	// 添加新端口
	for _, port := range newPorts {
		if err := tx.Create(&port).Error; err != nil {
			tx.Rollback()
			res.FailWithMsg("更新端口信息失败", c)
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		res.FailWithMsg("更新端口信息失败", c)
		return
	}

	msg := fmt.Sprintf("新增端口%d个，删除端口%d个", len(newPorts), len(portsToDelete))
	var portList []models.HoneyPortModel
	err = global.DB.Find(&portList, "honey_ip_id = ?", cr.HoneyIPID).Error

	req := mq_service.BindPortRequest{
		IP:    honeyIPModel.IP,
		LogID: "",
	}

	for _, model := range portList {
		req.PortList = append(req.PortList, mq_service.PortInfo{
			IP:       honeyIPModel.IP,
			Port:     model.Port,
			DestIP:   model.DstIP,
			DestPort: model.DstPort,
		})
	}
	mq_service.SendBindPortMsg(nodeModel.Uid, req)
	res.OkWithMsg(msg, c)
}
