package metrics

type HostMetrics struct {
	CpuUsage     float32 `json:"cpu_usage"`
	RamAvailable uint64  `json:"ram_available"`
	RamTotal     uint64  `json:"ram_total"`
}
