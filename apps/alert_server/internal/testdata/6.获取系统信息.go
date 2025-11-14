package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SystemInfo struct {
	Distribution  string    // 如 "CentOS 7", "Ubuntu 22.04"
	KernelVersion string    // 如 "5.15.0-86-generic"
	Arch          string    // "x86", "ARM", "unknown"
	BootTime      time.Time // 系统启动时间（本地时区）
}

// getDistribution 获取 Linux 发行版信息
func getDistribution() string {
	// 尝试读取 /etc/os-release（现代 Linux 标准）
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		var name, version string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "NAME=") {
				name = strings.Trim(strings.Split(line, "=")[1], `"`)
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				version = strings.Trim(strings.Split(line, "=")[1], `"`)
			}
		}
		if name != "" {
			if version != "" {
				return name + " " + version
			}
			return name
		}
	}

	// 备用：读取 /etc/issue（旧系统）
	if data, err := os.ReadFile("/etc/issue"); err == nil {
		issue := strings.TrimSpace(string(data))
		// 简单提取，例如 "CentOS Linux 7 \n \l" → 提取 "CentOS Linux 7"
		if parts := strings.Fields(issue); len(parts) >= 3 {
			return strings.Join(parts[:3], " ")
		}
		return issue
	}

	return "unknown"
}

// getKernelVersion 获取内核版本
func getKernelVersion() string {
	// 方法1：通过 /proc/version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) >= 3 {
			return parts[2] // 通常是第三个字段
		}
	}

	// 方法2：调用 uname -r
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}

	return "unknown"
}

// getArch 获取系统架构并映射为通用类型
func getArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64", "386", "x86_64":
		return "x86"
	case "arm", "arm64", "aarch64":
		return "ARM"
	case "ppc64", "ppc64le":
		return "PowerPC"
	case "s390x":
		return "IBM Z"
	default:
		return arch // 或 "unknown"
	}
}

// getBootTime 获取系统启动时间
func getBootTime() (time.Time, error) {
	// 读取 /proc/uptime，第一列为系统运行秒数
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 1 {
			return time.Time{}, fmt.Errorf("invalid /proc/uptime format")
		}
		uptimeSec, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return time.Time{}, err
		}
		now := time.Now()
		bootTime := now.Add(-time.Duration(uptimeSec*1e9) * time.Nanosecond)
		return bootTime, nil
	}
	return time.Time{}, fmt.Errorf("failed to read /proc/uptime")
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (*SystemInfo, error) {
	dist := getDistribution()
	kernel := getKernelVersion()
	arch := getArch()

	bootTime, err := getBootTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get boot time: %w", err)
	}

	return &SystemInfo{
		Distribution:  dist,
		KernelVersion: kernel,
		Arch:          arch,
		BootTime:      bootTime,
	}, nil
}

// 示例使用
func main() {
	info, err := GetSystemInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("发行版本: %s\n", info.Distribution)
	fmt.Printf("内核版本: %s\n", info.KernelVersion)
	fmt.Printf("系统架构: %s\n", info.Arch)
	fmt.Printf("启动时间: %s (本地时间)\n", info.BootTime.Format(time.DateTime))
}
