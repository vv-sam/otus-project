package task

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	STATUS_QUEUED      = 0
	STATUS_IN_PROGRESS = 1
	STATUS_OK          = 2
	STATUS_DELETED     = 3
)

type Task struct {
	// ID задачи
	Id uuid.UUID `json:"id" bson:"id"`

	// Статус задачи
	Status int16 `json:"status" bson:"status"`

	// Тип задачи
	Type string `json:"type" bson:"type"`
}

func (t Task) String() string {
	return fmt.Sprintf("%q, %d, %q", t.Id, t.Status, t.Type)
}

func (t Task) GetId() uuid.UUID {
	return t.Id
}

func (t Task) Validate() error {
	if t.Id == uuid.Nil {
		return fmt.Errorf("id is required")
	}

	if t.Type == "" {
		return fmt.Errorf("type is required")
	}

	return nil
}
