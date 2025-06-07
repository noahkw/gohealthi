package healthstats

import (
	"fmt"
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

	fmt.Println(memStats.UsedPercent)

	return memStats.UsedPercent, nil
}

func CpuUsage() ([]float64, error) {
	cpuPercent, err := cpu.Percent(time.Second, true)

	if err != nil {
		return []float64{}, err
	}

	fmt.Println(cpuPercent)

	return cpuPercent, nil
}

func DiskUsage() (float64, error) {
	diskStats, err := disk.Usage("/")

	if err != nil {
		return 0, err
	}

	fmt.Println(diskStats.UsedPercent)

	return diskStats.UsedPercent, nil
}

func NetworkUsage() (uint64, uint64, error) {
	netStats, err := net.IOCounters(false)

	if err != nil {
		return 0, 0, err
	}

	fmt.Println("recv", netStats[0].BytesRecv/MEGA)
	fmt.Println("sent", netStats[0].BytesSent/MEGA)

	return netStats[0].BytesRecv / MEGA, netStats[0].BytesSent / MEGA, nil
}
