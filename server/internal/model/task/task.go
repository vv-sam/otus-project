package task

import "github.com/google/uuid"

const (
	STATUS_QUEUED      = 0
	STATUS_IN_PROGRESS = 1
	STATUS_OK          = 2
	STATUS_DELETED     = 3
)

type Task struct {
	// ID задачи
	Id uuid.UUID `json:"id"`

	// Статус задачи
	Status int16 `json:"status"`

	// Тип задачи
	Type string `json:"type"`
}
