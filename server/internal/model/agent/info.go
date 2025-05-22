package agent

import (
	"fmt"

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

// Вернём строку с id агента и id статуса
func (i Info) String() string {
	return fmt.Sprintf("%q, %d", i.AgentId.String(), i.Status)
}

func (i Info) GetId() uuid.UUID {
	return i.AgentId
}

func (i Info) Validate() error {
	if i.AgentId == uuid.Nil {
		return fmt.Errorf("agent_id is required")
	}

	return nil
}
