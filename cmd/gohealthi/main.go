package main

import (
	"flag"
	"fmt"
	"github.com/noahkw/gohealthi/pkg/healthstats"
	"github.com/noahkw/gohealthi/pkg/models"
	"github.com/noahkw/gohealthi/pkg/server"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %v\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	fmt.Println("Printing some health stats")
	systemUsage, err := healthstats.CurrentSystemUsage()

	if err != nil {
		panic(err)
	}

	formatSystemUsage(*systemUsage)

	server.Serve()
}

func formatSystemUsage(usage models.SystemUsage) {
	fmt.Printf("RAM Usage: %.2f%%\n", usage.RamUsage)
	fmt.Printf("Disk Usage: %.2f%%\n", usage.DiskUsage)
	fmt.Println("CPU Percentages:")
	for i, cpuPercent := range usage.CpuPercentage {
		fmt.Printf("  CPU %d: %.2f%%\n", i, cpuPercent)
	}
	fmt.Printf("Average CPU Usage: %.2f%%\n", usage.CpuPercentageAvg)
	fmt.Printf("Network Bytes Received: %d MB\n", usage.NetworkBytesRecv)
	fmt.Printf("Network Bytes Sent: %d MB\n", usage.NetworkBytesSent)
}
