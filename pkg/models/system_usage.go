package models

import health "github.com/noahkw/gohealthi/proto"

type SystemUsage struct {
	RamUsage         float64
	DiskUsage        float64
	CpuPercentage    []float64
	CpuPercentageAvg float64
	NetworkBytesRecv uint64
	NetworkBytesSent uint64
}

func ToHealthResponse(systemUsage *SystemUsage) *health.HealthResponse {
	return &health.HealthResponse{
		RamUsage:         systemUsage.RamUsage,
		DiskUsage:        systemUsage.DiskUsage,
		CpuPercentage:    systemUsage.CpuPercentage,
		CpuPercentageAvg: systemUsage.CpuPercentageAvg,
		NetworkBytesRecv: systemUsage.NetworkBytesRecv,
		NetworkBytesSent: systemUsage.NetworkBytesSent,
	}
}

func NewSystemUsage(ramUsage, diskUsage float64, cpuPercentage []float64, cpuAvg float64, recv, sent uint64) *SystemUsage {
	return &SystemUsage{
		RamUsage:         ramUsage,
		DiskUsage:        diskUsage,
		CpuPercentage:    cpuPercentage,
		CpuPercentageAvg: cpuAvg,
		NetworkBytesRecv: recv,
		NetworkBytesSent: sent,
	}
}
