package healthstats

import (
	"fmt"
	"time"

	health "github.com/noahkw/gohealthi/proto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

const MEGA = 1_000_000

func RamUsage() (float64, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return memStats.UsedPercent, nil
}

func CpuUsage() ([]float64, error) {
	cpuPercent, err := cpu.Percent(time.Second, true)

	if err != nil {
		return []float64{}, err
	}

	return cpuPercent, nil
}

func DiskUsage() (float64, error) {
	diskStats, err := disk.Usage("/")

	if err != nil {
		return 0, err
	}

	return diskStats.UsedPercent, nil
}

func NetworkUsage() (uint64, uint64, error) {
	netStats, err := net.IOCounters(false)

	if err != nil {
		return 0, 0, err
	}

	return netStats[0].BytesRecv / MEGA, netStats[0].BytesSent / MEGA, nil
}

func SystemUsageMean(healthResponses []health.HealthResponse) (*health.HealthResponse, error) {
	if len(healthResponses) == 0 {
		return nil, fmt.Errorf("no health responses provided")
	}

	var totalRamUsage, totalDiskUsage float64
	var totalNetworkBytesRecv, totalNetworkBytesSent uint64
	var totalCpuPercentageAvg float64

	for _, response := range healthResponses {
		totalRamUsage += response.RamUsage
		totalDiskUsage += response.DiskUsage
		totalNetworkBytesRecv += response.NetworkBytesRecv
		totalNetworkBytesSent += response.NetworkBytesSent

		totalCpuPercentageAvg += response.CpuPercentageAvg
	}

	count := float64(len(healthResponses))

	meanRamUsage := totalRamUsage / count
	meanDiskUsage := totalDiskUsage / count
	meanNetworkBytesRecv := totalNetworkBytesRecv / uint64(len(healthResponses))
	meanNetworkBytesSent := totalNetworkBytesSent / uint64(len(healthResponses))

	return &health.HealthResponse{
		RamUsage:         meanRamUsage,
		DiskUsage:        meanDiskUsage,
		CpuPercentage:    make([]float64, 0),
		CpuPercentageAvg: totalCpuPercentageAvg / count,
		NetworkBytesRecv: meanNetworkBytesRecv,
		NetworkBytesSent: meanNetworkBytesSent,
	}, nil
}

func CurrentSystemUsage() (health.HealthResponse, error) {
	ramUsage, err := RamUsage()
	if err != nil {
		return health.HealthResponse{}, err
	}

	diskUsage, err := DiskUsage()
	if err != nil {
		return health.HealthResponse{}, err
	}

	cpuPercentages, err := CpuUsage()
	if err != nil {
		return health.HealthResponse{}, err
	}

	networkBytesRecv, networkBytesSent, err := NetworkUsage()
	if err != nil {
		return health.HealthResponse{}, err
	}

	return health.HealthResponse{
		RamUsage:         ramUsage,
		DiskUsage:        diskUsage,
		CpuPercentage:    cpuPercentages,
		CpuPercentageAvg: averageCpuPercentage(cpuPercentages),
		NetworkBytesRecv: networkBytesRecv,
		NetworkBytesSent: networkBytesSent,
	}, nil
}

func averageCpuPercentage(percentages []float64) float64 {
	if len(percentages) == 0 {
		return 0
	}

	var total float64
	for _, percentage := range percentages {
		total += percentage
	}

	return total / float64(len(percentages))
}
