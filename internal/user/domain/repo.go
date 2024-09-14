package domain

import (
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, u *User) (*User, error)
}
