package metrics

import "fmt"

type HostMetrics struct {
	CpuUsage     float32 `json:"cpu_usage" bson:"cpu_usage"`
	RamAvailable uint64  `json:"ram_available" bson:"ram_available"`
	RamTotal     uint64  `json:"ram_total" bson:"ram_total"`
}

func (m HostMetrics) String() string {
	return fmt.Sprintf("%f, %d / %d", m.CpuUsage, m.RamAvailable, m.RamTotal)
}
