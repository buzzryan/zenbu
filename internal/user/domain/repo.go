package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, u *User) (*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
}
