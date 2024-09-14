package infra

import (
	"context"

	"github.com/buzzryan/zenbu/internal/user/domain"
)

type userRepo struct{}

func NewUserRepo() domain.UserRepo {
	return &userRepo{}
}

func (ur *userRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	return nil, nil
}
