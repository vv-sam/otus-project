package service

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/agent"
	"github.com/vv-sam/otus-project/server/internal/model/configuration"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
	"github.com/vv-sam/otus-project/server/internal/repository"
)

func GenerateStruct(ch chan fmt.Stringer) {
	// рандомная моделька
	t := rand.Uint32() % 4

	switch t {
	case 0:
		ch <- &agent.Info{AgentId: uuid.New(), Status: 0}
	case 1:
		ch <- &configuration.Minecraft{ServerName: "Otus server", MaxPlayers: 64}
	case 2:
		ch <- &metrics.HostMetrics{CpuUsage: float32(t) / 8, RamAvailable: 1_000_000, RamTotal: 3_000_000}
	case 3:
		ch <- &task.Task{Id: uuid.New(), Status: task.STATUS_QUEUED, Type: "test"}
	}
}

func ConsumeStructs(ch chan fmt.Stringer) {
	for s := range ch {
		repository.PassStruct(s)
	}
}
