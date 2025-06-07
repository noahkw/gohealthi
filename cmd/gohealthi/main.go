package main

import (
	"fmt"
	"github.com/noahkw/gohealthi/pkg/healthstats"
)

func main() {
	fmt.Println("Printing some health stats")

	healthstats.RamUsage()
	healthstats.CpuUsage()
	healthstats.DiskUsage()
	healthstats.NetworkUsage()
}
