package core

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

//go:embed oui.txt
var oui []byte // 编译时嵌入 oui.txt 文件内容

var ouiMap map[string]string

// loadOUI 从内置的 oui 数据加载（格式：前缀 + 厂商名）
func loadOUI() (map[string]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(oui)))
	ouiMap := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		prefix := strings.ToUpper(strings.ReplaceAll(fields[0], "-", ":"))
		ouiMap[prefix] = strings.Join(fields[1:], " ")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 OUI 数据错误: %v", err)
	}

	return ouiMap, nil
}

// GetVendorByMAC 查询厂商
func GetVendorByMAC(mac string) string {
	mac = strings.ToUpper(mac)
	parts := strings.Split(mac, ":")
	if len(parts) < 3 {
		return "无效的 MAC 地址"
	}
	prefix := strings.Join(parts[:3], ":")
	if vendor, ok := ouiMap[prefix]; ok {
		return vendor
	}
	return "未知厂商"
}

func init() {
	var err error
	ouiMap, err = loadOUI()
	if err != nil {
		panic(err)
	}
	logrus.Infof("加载 OUI 数据成功，共有 %d 条数据", len(ouiMap))
}
