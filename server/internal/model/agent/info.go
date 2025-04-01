package agent

import (
	"github.com/google/uuid"
	"github.com/vv-sam/otus-project/server/internal/model/metrics"
	"github.com/vv-sam/otus-project/server/internal/model/task"
)

type Info struct {
	AgentId      uuid.UUID           `json:"agent_id"`
	Status       int16               `json:"status"`
	CurrentTasks []task.Task         `json:"tasks"`
	Metrics      metrics.HostMetrics `json:"metrics"`
}
