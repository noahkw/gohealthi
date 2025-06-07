package main

import (
	"fmt"
	"github.com/noahkw/gohealthi/pkg/healthstats"
)

func main() {
	fmt.Println("Printing some health stats")

	systemUsage, err := healthstats.CurrentSystemUsage()

	if err != nil {
		panic(err)
	}

	formatSystemUsage(systemUsage)
}

func formatSystemUsage(usage healthstats.SystemUsage) {
	fmt.Printf("RAM Usage: %.2f%%\n", usage.RamUsage)
	fmt.Printf("Disk Usage: %.2f%%\n", usage.DiskUsage)
	fmt.Println("CPU Percentages:")
	for i, cpuPercent := range usage.CpuPercentages {
		fmt.Printf("  CPU %d: %.2f%%\n", i, cpuPercent)
	}
	fmt.Printf("Network Bytes Received: %d MB\n", usage.NetworkBytesRecv)
	fmt.Printf("Network Bytes Sent: %d MB\n", usage.NetworkBytesSent)
}
