package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID
	ShortID string
	Server  string
	Tag     string
	Email   string
}
