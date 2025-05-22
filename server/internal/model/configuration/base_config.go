package configuration

import (
	"fmt"

	"github.com/google/uuid"
)

// Базовые опции, присущие каждой возможной конфигурации
type baseConfig struct {
	// уникальный ID конфигурации
	Id uuid.UUID

	// ID агента для запуска задачи
	AgentId uuid.UUID

	// Порт сервера
	Port uint16 `json:"port"`

	// Тип конфигурации
	Type string `json:"type"`
}

func (c *baseConfig) GetId() uuid.UUID {
	return c.Id
}

func (c *baseConfig) Validate() error {
	if c.Id == uuid.Nil {
		return fmt.Errorf("id is required")
	}

	if c.Type == "" {
		return fmt.Errorf("type is required")
	}

	return nil
}
