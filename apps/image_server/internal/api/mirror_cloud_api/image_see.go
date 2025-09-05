package mirror_cloud_api

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"image_server/internal/utils/res"
)

type ImageSeeResponse struct {
	ImageID   string `json:"imageID"`
	ImageName string `json:"imageName"` // Repository
	ImageTag  string `json:"imageTag"`
	ImagePath string `json:"imagePath"`
}

// ImageSeeView 处理镜像上传并解析基本信息
func (MirrorCloudApi) ImageSeeView(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		res.FailWithMsg("请选择镜像文件", c)
		return
	}

	// 支持 .tar 和 .tar.gz
	filename := file.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	baseName := filename
	var isGzipped bool

	switch ext {
	case ".tar":
		isGzipped = false
	case ".gz":
		// 检查是否是 .tar.gz
		if strings.HasSuffix(strings.ToLower(filename), ".tar.gz") {
			isGzipped = true
		} else {
			res.FailWithMsg("不支持的文件格式，请上传 .tar 或 .tar.gz 文件", c)
			return
		}
	default:
		res.FailWithMsg("不支持的文件格式，请上传 .tar 或 .tar.gz 文件", c)
		return
	}

	logrus.Infof("ext is %v, baseName is %v", ext, baseName)

	// 创建临时文件名避免冲突
	tempFileName := baseName

	// 路径
	tempDir := "uploads/images_temp"
	finalDir := "uploads/images"
	tempPath := filepath.Join(tempDir, tempFileName)

	// 确保目录存在
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		res.FailWithMsg("无法创建临时目录", c)
		return
	}
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		res.FailWithMsg("无法创建目标目录", c)
		return
	}

	// 保存上传的文件
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		res.FailWithMsg("文件保存失败", c)
		return
	}

	// 解析镜像信息
	imageID, imageName, imageTag, err := parseDockerTar(tempPath, isGzipped)
	if err != nil {
		_ = os.Remove(tempPath) // 解析失败，清理临时文件
		res.FailWithMsg("无法解析镜像文件: "+err.Error(), c)
		return
	}

	// 构造响应
	data := ImageSeeResponse{
		ImageID:   imageID,
		ImageName: imageName,
		ImageTag:  imageTag,
		ImagePath: "/uploads/images/" + filename, // 前端可用路径（配合静态文件服务）
	}

	// 异步任务：5分钟后移动文件
	go func() {
		time.Sleep(10 * time.Second)
		finalPath := filepath.Join(finalDir, filename)

		// 检查目标是否存在
		if _, err := os.Stat(finalPath); err == nil {
			// 已存在，重命名避免覆盖
			finalPath = filepath.Join(finalDir, filename)
		}

		if err := os.Rename(tempPath, finalPath); err != nil {
			// 移动失败，可记录日志或重试
			fmt.Printf("移动镜像文件失败: %v\n", err)
		} else {
			fmt.Printf("镜像文件已移动: %s -> %s\n", tempPath, finalPath)
		}
	}()

	res.OkWithData(data, c)
}

// parseDockerTar 解析 Docker 镜像 tar 包，提取 Image ID、Repository、Tag
func parseDockerTar(filePath string, isGzipped bool) (imageID, imageName, imageTag string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()

	var reader io.Reader = f

	// 如果是 .tar.gz，先解压
	if isGzipped {
		gzReader, err := gzip.NewReader(f)
		if err != nil {
			return "", "", "", fmt.Errorf("解压失败: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tarReader := tar.NewReader(reader)

	var manifestData []byte
	var repositoriesData []byte

	// 遍历 tar 文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", "", err
		}

		switch header.Name {
		case "manifest.json":
			manifestData, err = io.ReadAll(tarReader)
			if err != nil {
				return "", "", "", err
			}
		case "repositories":
			repositoriesData, err = io.ReadAll(tarReader)
			if err != nil {
				return "", "", "", err
			}
		}
	}

	// 优先解析 manifest.json（标准方式）
	if len(manifestData) > 0 {
		var manifest []map[string]interface{}
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			return "", "", "", err
		}
		if len(manifest) > 0 {
			configFile := manifest[0]["Config"].(string)
			imageID = strings.Split(configFile, "/")[2][:12]

			// Repositories
			repos, ok := manifest[0]["RepoTags"].([]interface{})
			if ok && len(repos) > 0 {
				tag := repos[0].(string) // 格式: repo:tag
				parts := strings.Split(tag, ":")
				imageName = strings.Join(parts[:len(parts)-1], ":") // 兼容含 : 的 repo 名
				imageTag = parts[len(parts)-1]
			} else {
				imageName = "<unknown>"
				imageTag = "<unknown>"
			}
			return imageID, imageName, imageTag, nil
		}
	}

	// fallback: 解析 repositories 文件（旧版 Docker）
	if len(repositoriesData) > 0 {
		var repos map[string]map[string]string
		if err := json.Unmarshal(repositoriesData, &repos); err != nil {
			return "", "", "", err
		}
		for repo, tags := range repos {
			for tag, id := range tags {
				imageID = id
				imageName = repo
				imageTag = tag
				return imageID, imageName, imageTag, nil
			}
		}
	}

	return "", "", "", fmt.Errorf("无法解析镜像元信息")
}
