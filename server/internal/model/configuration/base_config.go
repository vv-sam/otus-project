package configuration

import "github.com/google/uuid"

// Базовые опции, присущие каждой возможной конфигурации
type baseConfig struct {
	// ID агента для запуска задачи
	AgentId uuid.UUID

	// Порт сервера
	Port uint16 `json:"port"`

	// Тип конфигурации
	Type string `json:"type"`
}
