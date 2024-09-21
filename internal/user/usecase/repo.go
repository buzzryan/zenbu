package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/user/domain"
)

// UserRepo is the interface that wraps the basic CRUD operations for User entity. (port)
type UserRepo interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByName(ctx context.Context, name string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) (*domain.User, error)
}
