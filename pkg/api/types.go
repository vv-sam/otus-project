package api

import "github.com/google/uuid"

type TaskStatus uint8

const (
	TASK_STATUS_QUEUED      TaskStatus = 0
	TASK_STATUS_IN_PROGRESS TaskStatus = 1
	TASK_STATUS_EXECUTING   TaskStatus = 2
	TASK_STATUS_DONE        TaskStatus = 3
	TASK_STATUS_ERROR       TaskStatus = 4
)

type TaskType uint8

const (
	TASK_TYPE_MINECRAFT TaskType = 1
	TASK_TYPE_FACTORIO  TaskType = 2
)

// ещё нужно продумать хоть немного безопасный механизм регистрации и авторизации агентов
type RegisterAgentRequest struct {
	Token string `json:"token"`
}

type RegisterAgentResponse struct {
}

type ConnectAgentRequest struct {
}

type ConnectAgentResponse struct {
}

// тоже нужно ещё подумать
type SendTaskRequest struct {
	Type TaskType `json:"type"`

	// К каждому типу таски нужно передавать, по идее, уникальный конфиг, т.к. для разных игр могут понадобиться разные настройки
	Options []byte `json:"options"`
}

// Агент вернёт id задачи, по которому можно будет отслеживать статус и отфильтровывать события
type SendTaskResponse struct {
	TaskId uuid.UUID `json:"task_id"`
}

type TaskStatusRequest struct {
	TaskId uuid.UUID `json:"task_id"`
}

// Нужно подумать насчёт логов, как передавать их с агента на сервер и отображать, и стоит ли вообще
type TaskStatusResponse struct {
	TaskId uuid.UUID  `json:"task_id"`
	Status TaskStatus `json:"status"`
}
