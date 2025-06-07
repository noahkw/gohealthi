package healthstats

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
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

type SystemUsage struct {
	RamUsage         float64
	DiskUsage        float64
	CpuPercentages   []float64
	NetworkBytesRecv uint64
	NetworkBytesSent uint64
}

func CurrentSystemUsage() (SystemUsage, error) {
	ramUsage, err := RamUsage()
	if err != nil {
		return SystemUsage{}, err
	}

	diskUsage, err := DiskUsage()
	if err != nil {
		return SystemUsage{}, err
	}

	cpuPercentages, err := CpuUsage()
	if err != nil {
		return SystemUsage{}, err
	}

	networkBytesRecv, networkBytesSent, err := NetworkUsage()
	if err != nil {
		return SystemUsage{}, err
	}

	return SystemUsage{
		RamUsage:         ramUsage,
		DiskUsage:        diskUsage,
		CpuPercentages:   cpuPercentages,
		NetworkBytesRecv: networkBytesRecv,
		NetworkBytesSent: networkBytesSent,
	}, nil
}
