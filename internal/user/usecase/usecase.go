package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/user/domain"
)

type SignupReq struct {
	Username string
	Password string
}

type SignupRes struct {
	Token string
}

type BasicSignupUC interface {
	Execute(ctx context.Context, req *SignupReq) (*SignupRes, error)
}

type basicSignupUC struct {
	userRepo     domain.UserRepo
	tokenManager TokenManager
}

func NewBasicSignupUC(userRepo domain.UserRepo, manager TokenManager) BasicSignupUC {
	return &basicSignupUC{userRepo: userRepo, tokenManager: manager}
}

func (b *basicSignupUC) Execute(ctx context.Context, req *SignupReq) (*SignupRes, error) {
	newUser := &domain.User{ID: uuid.New(), Username: req.Username, Password: domain.NewPassword(req.Password)}
	newUser, err := b.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	token, err := b.tokenManager.Generate(&Claims{
		UserID: newUser.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SignupRes{Token: token}, nil
}
