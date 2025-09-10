package mirror_cloud_api

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"
)

type ImageCreateRequest struct {
	ImageID   string `json:"imageID" binding:"required"`
	ImageName string `json:"imageName" binding:"required"` // Repository
	ImageTag  string `json:"imageTag" binding:"required"`
	ImagePath string `json:"imagePath" binding:"required"`            // 镜像上传的路径
	Title     string `json:"title" binding:"required"`                // 镜像别名
	Port      int    `json:"port" binding:"required,min=1,max=65535"` // 镜像端口
	Agreement int8   `json:"agreement" binding:"required,oneof=1"`    // 镜像协议
}

func (MirrorCloudApi) ImageCreateView(c *gin.Context) {
	cr := middleware.GetBind[ImageCreateRequest](c)

	// 1. 校验 title 是否重复
	var model models.ImageModel
	if err := global.DB.Take(&model, "title = ?", cr.Title).Error; err == nil {
		res.FailWithMsg("镜像别名不能重复", c)
		return
	}

	// 2. 校验 imageName + imageTag 是否重复
	if err := global.DB.Take(&model, "image_name = ? AND tag = ?", cr.ImageName, cr.ImageTag).Error; err == nil {
		res.FailWithMsg("镜像名+标签已存在，不能重复导入", c)
		return
	}

	basePath, _ := os.Getwd()
	fullPath := filepath.Join(basePath, cr.ImagePath)
	// 3. 校验镜像文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		res.FailWithMsg("镜像文件不存在，请重新上传", c)
		return
	}

	// 4. 执行 docker load 导入镜像
	cmd := exec.Command("docker", "load", "-i", fullPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.Log.Error("docker load 失败:", string(output))
		res.FailWithMsg("镜像导入失败，请检查镜像格式或Docker服务状态", c)
		return
	}

	// 5. 移动镜像文件到正式目录
	//targetDir := "uploads/images"
	//targetPath := filepath.Join(targetDir, filepath.Base(cr.ImagePath))

	//if err := os.MkdirAll(targetDir, 0755); err != nil {
	//	global.Log.Error("创建目标目录失败:", err)
	//	res.FailWithMsg("系统错误：无法创建存储目录", c)
	//	return
	//}
	//
	//if err := os.Rename(cr.ImagePath, targetPath); err != nil {
	//	global.Log.Error("移动镜像文件失败:", err)
	//	res.FailWithMsg("系统错误：移动镜像失败", c)
	//	return
	//}

	// 6. 数据入库
	image := models.ImageModel{
		DockerImageID: cr.ImageID,
		ImageName:     cr.ImageName,
		Tag:           cr.ImageTag,
		ImagePath:     cr.ImagePath, // 存储移动后的新路径
		Title:         cr.Title,
		Port:          cr.Port,
		Agreement:     cr.Agreement,
	}

	if err := global.DB.Create(&image).Error; err != nil {
		global.Log.Error("保存镜像记录失败:", err)
		res.FailWithMsg("系统错误：保存数据失败", c)
		return
	}

	//  成功响应
	res.OkWithMsg("镜像创建成功", c)
}
