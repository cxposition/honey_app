package net_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
)

type OptionsResponse struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

func (NetApi) OptionsView(c *gin.Context) {
	var netList []models.NetModel
	global.DB.Find(&netList)
	var list = make([]OptionsResponse, 0)
	for _, model := range netList {
		list = append(list, OptionsResponse{
			Label: fmt.Sprintf("%s(%s)", model.Title, model.Subnet()),
			Value: model.ID,
		})
	}

	res.OkWithData(list, c)
}
