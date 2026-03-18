package repository

import (
	"context"

	"github.com/VAGRAMCHIC/vless_reality_agent/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetAll(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id *uuid.UUID) (*domain.User, error)
	DeleteUser(ctx context.Context, id *uuid.UUID) error
}
