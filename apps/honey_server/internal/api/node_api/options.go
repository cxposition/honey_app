package node_api

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

func (NodeApi) OptionsView(c *gin.Context) {
	var nodeList []models.NodeModel
	global.DB.Find(&nodeList)
	var list = make([]OptionsResponse, 0)
	for _, model := range nodeList {
		list = append(list, OptionsResponse{
			Value: model.ID,
			Label: fmt.Sprintf("%s(%s)", model.Title, model.IP),
		})
	}

	res.OkWithData(list, c)
}
