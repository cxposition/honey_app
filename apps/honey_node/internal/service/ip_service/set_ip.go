package ip_service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"honey_node/internal/utils/cmd"
	"strings"
)

type SetIpRequest struct {
	Ip       string `json:"ip"`
	Mask     int8   `json:"mask"`
	LinkName string `json:"linkName"` // 自己的接口名称
	Network  string `json:"network"`  // 基于哪个网卡创建的
	Mac      string `json:"mac"`
}

func SetIp(req SetIpRequest) (mac string, err error) {
	linkName := req.LinkName
	// 创建资源清理函数
	cleanup := func() {
		// 失败时尝试清理已创建的资源
		if err := cmd.Cmd(fmt.Sprintf("ip link delete %s", linkName)); err != nil {
			logrus.Errorf("清理失败，删除网络接口 %s 时出错: %v", linkName, err)
		}
	}

	// 执行网络配置命令
	if err = createMacVlanInterface(linkName, req.Network); err != nil {
		logrus.Errorf("创建macvlan接口失败: %v", err)
		cleanup()
		return
	}

	if req.Mac != "" {
		err = setInterfaceMac(linkName, req.Mac)
		if err != nil {
			logrus.Errorf("设置mac失败: %v", err)
			cleanup()
			return
		}
	}

	if err = setInterfaceUp(linkName); err != nil {
		logrus.Errorf("启用网络接口失败: %v", err)
		cleanup()
		return
	}

	if err = addIPAddress(linkName, req.Ip, req.Mask); err != nil {
		logrus.Errorf("添加IP地址失败: %v", err)
		cleanup()
		return
	}

	if req.Mac == "" {
		req.Mac, err = GetMACAddress(linkName)
		if err != nil {
			return
		}
	}
	return req.Mac, nil
}

// 创建macvlan接口
func createMacVlanInterface(linkName, network string) error {
	cmdStr := fmt.Sprintf("ip link add %s link %s type macvlan mode bridge", linkName, network)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// 启用网络接口
func setInterfaceUp(linkName string) error {
	cmdStr := fmt.Sprintf("ip link set %s up", linkName)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// 设置mac地址
func setInterfaceMac(linkName string, mac string) error {
	cmdStr := fmt.Sprintf("ip link set %s address %s", linkName, mac)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// addIPAddress 添加IP地址
func addIPAddress(linkName, ip string, mask int8) error {
	cmdStr := fmt.Sprintf("ip addr add %s/%d dev %s", ip, mask, linkName)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// GetMACAddress 获取接口MAC地址
func GetMACAddress(linkName string) (string, error) {
	cmdStr := fmt.Sprintf("ip link show %s | awk '/link\\/ether/ {print $2}'", linkName)
	mac, err := cmd.CommandWithOut(cmdStr)
	if err != nil {
		return "", fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return strings.TrimSpace(mac), nil
}

func RemoveInterface(iface string) error {
	err := cmd.Cmd(fmt.Sprintf("ip link del %s", iface))
	return err
}
