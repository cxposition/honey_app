package vs_net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/config"
	"image_server/internal/core"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/cmd"
	"image_server/internal/utils/res"
)

type VsNetApi struct {
}

type VsNetRequest struct {
	Name   string `json:"name" binding:"required"`
	Prefix string `json:"prefix" binding:"required"`
	Net    string `json:"net" binding:"required"`
}

func (VsNetApi) VsNetInfoView(c *gin.Context) {
	res.OkWithData(global.Config.VsNet, c)
}

func (VsNetApi) VsNetUpdateView(c *gin.Context) {
	cr := middleware.GetBind[VsNetRequest](c)
	// 在没有虚拟服务的情况下才能创建
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList)
	if len(serviceList) != 0 {
		res.FailWithMsg("存在虚拟服务，不可修改虚拟子网", c)
		return
	}

	// 把之前的删掉
	command := fmt.Sprintf("docker network rm %s", global.Config.VsNet.Name)
	err := cmd.Cmd(command)
	if err != nil {
		res.FailWithMsg("删除虚拟子网失败", c)
		return
	}

	// 创建新的
	command = fmt.Sprintf("docker network create --driver bridge --subnet %s %s", cr.Net, cr.Name)
	err = cmd.Cmd(command)
	if err != nil {
		res.FailWithMsg("创建虚拟子网失败", c)
		return
	}

	// 回写到配置文件中
	global.Config.VsNet = config.VsNet{
		Name:   cr.Name,
		Prefix: cr.Prefix,
		Net:    cr.Net,
	}

	core.SetConfig()
	res.OkWithMsg("修改虚拟网络成功", c)
}
