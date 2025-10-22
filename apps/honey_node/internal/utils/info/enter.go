package info

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type ResourceMessage struct {
	CPUCount              int64
	CPUUseRate            float64 // %
	MemTotal              int64   // GB
	MemUseRate            float64 // %
	DiskTotal             int64   // GB
	DiskUseRate           float64 // %
	NodePath              string
	NodeResourceOccupancy int64 // %
}

// 字节转 GB（四舍五入）
func bytesToGB(b uint64) int64 {
	const gb = 1024 * 1024 * 1024
	return int64(math.Round(float64(b) / float64(gb)))
}

// 四舍五入保留两位小数
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// GetSystemResource 获取系统资源信息（单位已转换）
func GetSystemResource() (*ResourceMessage, error) {
	// CPU 信息
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("get cpu count failed: %v", err)
	}
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("get cpu percent failed: %v", err)
	}
	cpuUseRate := 0.0
	if len(cpuPercents) > 0 {
		cpuUseRate = round2(cpuPercents[0])
	}

	// 内存信息
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get memory info failed: %v", err)
	}
	memTotal := bytesToGB(vmStat.Total)
	memUseRate := round2(vmStat.UsedPercent)

	// 磁盘信息（根目录）
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("get disk info failed: %v", err)
	}
	diskTotal := bytesToGB(diskStat.Total)
	diskUseRate := round2(diskStat.UsedPercent)

	// 主机名
	hostname, _ := os.Hostname()

	// 综合占用率
	nodeResourceOccupancy := int64(math.Round((cpuUseRate + memUseRate + diskUseRate) / 3))

	return &ResourceMessage{
		CPUCount:              int64(cpuCount),
		CPUUseRate:            cpuUseRate,
		MemTotal:              memTotal,
		MemUseRate:            memUseRate,
		DiskTotal:             diskTotal,
		DiskUseRate:           diskUseRate,
		NodePath:              hostname,
		NodeResourceOccupancy: nodeResourceOccupancy,
	}, nil
}

type SystemInfo struct {
	Distribution  string    // 如 "CentOS 7", "Ubuntu 22.04"
	KernelVersion string    // 如 "5.15.0-86-generic"
	Arch          string    // "x86", "ARM", "unknown"
	BootTime      time.Time // 系统启动时间（本地时区）
}

// GetDistribution 获取 Linux 发行版信息
func GetDistribution() string {
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

// GetKernelVersion 获取内核版本
func GetKernelVersion() string {
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

// GetArch 获取系统架构并映射为通用类型
func GetArch() string {
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

// GetBootTime 获取系统启动时间
func GetBootTime() (time.Time, error) {
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
	dist := GetDistribution()
	kernel := GetKernelVersion()
	arch := GetArch()

	bootTime, err := GetBootTime()
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
