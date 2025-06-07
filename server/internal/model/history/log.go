package history

import (
	"time"

	"github.com/google/uuid"
)

type Log[T any] struct {
	Time   time.Time `json:"time"`
	Action string    `json:"action"`
	Id     uuid.UUID `json:"id"`
	Data   T         `json:"data"`
}
