package healthstats

import (
	"fmt"
	"github.com/noahkw/gohealthi/pkg/models"
	"time"

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

func SystemUsageMean(systemUsages []*models.SystemUsage) (*models.SystemUsage, error) {
	if len(systemUsages) == 0 {
		return nil, fmt.Errorf("no health responses provided")
	}

	var totalRamUsage, totalDiskUsage float64
	var totalNetworkBytesRecv, totalNetworkBytesSent uint64
	var totalCpuPercentageAvg float64

	for _, response := range systemUsages {
		totalRamUsage += response.RamUsage
		totalDiskUsage += response.DiskUsage
		totalNetworkBytesRecv += response.NetworkBytesRecv
		totalNetworkBytesSent += response.NetworkBytesSent

		totalCpuPercentageAvg += response.CpuPercentageAvg
	}

	count := float64(len(systemUsages))

	meanRamUsage := totalRamUsage / count
	meanDiskUsage := totalDiskUsage / count
	meanNetworkBytesRecv := totalNetworkBytesRecv / uint64(len(systemUsages))
	meanNetworkBytesSent := totalNetworkBytesSent / uint64(len(systemUsages))

	return models.NewSystemUsage(
		meanRamUsage,
		meanDiskUsage,
		[]float64{},
		totalCpuPercentageAvg/count,
		meanNetworkBytesRecv,
		meanNetworkBytesSent,
	), nil
}

func CurrentSystemUsage() (*models.SystemUsage, error) {
	ramUsage, err := RamUsage()
	if err != nil {
		return nil, err
	}

	diskUsage, err := DiskUsage()
	if err != nil {
		return nil, err
	}

	cpuPercentages, err := CpuUsage()
	if err != nil {
		return nil, err
	}

	networkBytesRecv, networkBytesSent, err := NetworkUsage()
	if err != nil {
		return nil, err
	}

	return models.NewSystemUsage(ramUsage, diskUsage, cpuPercentages, mean(cpuPercentages), networkBytesRecv, networkBytesSent), nil
}

func mean(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}

	var total float64
	for _, percentage := range arr {
		total += percentage
	}

	return total / float64(len(arr))
}
