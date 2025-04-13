package service

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
)

func GenerateStructs() []interface{} {
	// количество от 5 до 15
	c := rand.Uint32()%10 + 5

	r := make([]interface{}, 0, c)

	for range c {
		// рандомная моделька
		t := rand.Uint32() % 4

		switch t {
		case 0:
			r = append(r, &agent.Info{AgentId: uuid.New(), Status: 0})
		case 1:
			r = append(r, &configuration.Minecraft{ServerName: "Otus server", MaxPlayers: uint(c)})
		case 2:
			r = append(r, &metrics.HostMetrics{CpuUsage: float32(t) / 8, RamAvailable: 1_000_000, RamTotal: 3_000_000})
		case 3:
			r = append(r, &task.Task{Id: uuid.New(), Status: task.STATUS_QUEUED, Type: "test"})
		}
	}

	return r
}
