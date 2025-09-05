package mirror_cloud_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image_server/internal/utils/res"
)

func (MirrorCloudApi) ImageSeeView(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		res.FailWithMsg("请选择镜像文件", c)
		return
	}

	fmt.Println("file header:::::", file.Header)
	res.OkWithMsg("上传成功", c)
}
