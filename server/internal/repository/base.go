package repository

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("item not found")
)

type uniqueObject interface {
	GetId() uuid.UUID
}
