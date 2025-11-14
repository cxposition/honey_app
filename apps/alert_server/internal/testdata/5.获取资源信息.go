package main

import (
	"fmt"
	"math"
	"os"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// ResourceMessage 对应 protobuf 结构
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

func main() {
	res, err := GetSystemResource()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("系统资源信息：\n")
	fmt.Printf("CPU核数: %d\n", res.CPUCount)
	fmt.Printf("CPU使用率: %.2f%%\n", res.CPUUseRate)
	fmt.Printf("内存总量: %d GB\n", res.MemTotal)
	fmt.Printf("内存使用率: %.2f%%\n", res.MemUseRate)
	fmt.Printf("磁盘总量: %d GB\n", res.DiskTotal)
	fmt.Printf("磁盘使用率: %.2f%%\n", res.DiskUseRate)
	fmt.Printf("节点名: %s\n", res.NodePath)
	fmt.Printf("综合资源占用率: %d%%\n", res.NodeResourceOccupancy)
}
