package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// OUI 数据文件名
const OUIFile = "oui.txt"

// loadOUI 从当前目录加载 OUI 文件（格式：前缀 + 厂商名）
func loadOUI(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("无法打开 OUI 文件: %v", err)
	}
	defer f.Close()

	ouiMap := make(map[string]string)
	scanner := bufio.NewScanner(f)
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
		return nil, fmt.Errorf("读取 OUI 文件错误: %v", err)
	}
	return ouiMap, nil
}

// getVendorByMAC 查询厂商
func getVendorByMAC(mac string, ouiMap map[string]string) string {
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

func main() {
	ouiMap, err := loadOUI(OUIFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	macs := []string{
		"00:0C:29:C3:BB:6B", // VMware
		"FC:FB:FB:01:02:03", // Meta
		"F8:27:93:AA:BB:CC", // Apple
		"34:29:8F:11:22:33", // Xiaomi
		"D8:53:BC:33:DA:66",
	}

	for _, mac := range macs {
		fmt.Printf("MAC: %-20s 厂商: %s\n", mac, getVendorByMAC(mac, ouiMap))
	}
}
